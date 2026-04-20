package api

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/amir-sheikh/secscan/backend/internal/model"
)

func renderPDF(scan *model.Scan) []byte {
	lines := []string{
		"SecScan Report",
		fmt.Sprintf("Scan ID: %s", scan.ID),
		fmt.Sprintf("Target: %s", scan.URL),
		fmt.Sprintf("Grade: %s", scan.Summary.Grade),
		fmt.Sprintf("Score: %d", scan.Summary.Score),
		fmt.Sprintf("Risk Level: %s", strings.ToUpper(scan.Summary.RiskLevel)),
		fmt.Sprintf("Findings: %d", scan.Summary.Findings),
	}

	for name, result := range scan.Modules {
		lines = append(lines, fmt.Sprintf("%s -> %s (%d)", strings.ToUpper(name), result.Status, result.Score))
	}

	content := buildPDFText(lines)
	var out bytes.Buffer
	offsets := []int{0}
	writeObject := func(index int, body string) {
		offsets = append(offsets, out.Len())
		fmt.Fprintf(&out, "%d 0 obj\n%s\nendobj\n", index, body)
	}

	out.WriteString("%PDF-1.4\n")
	writeObject(1, "<< /Type /Catalog /Pages 2 0 R >>")
	writeObject(2, "<< /Type /Pages /Count 1 /Kids [3 0 R] >>")
	writeObject(3, "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>")
	writeObject(4, fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(content), content))
	writeObject(5, "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>")

	xrefPos := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n", len(offsets))
	out.WriteString("0000000000 65535 f \n")
	for _, offset := range offsets[1:] {
		fmt.Fprintf(&out, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&out, "trailer << /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF", len(offsets), xrefPos)
	return out.Bytes()
}

func buildPDFText(lines []string) string {
	escaped := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\\", "\\\\")
		line = strings.ReplaceAll(line, "(", "\\(")
		line = strings.ReplaceAll(line, ")", "\\)")
		escaped = append(escaped, fmt.Sprintf("(%s) Tj", line))
	}

	builder := &strings.Builder{}
	builder.WriteString("BT\n/F1 12 Tf\n50 780 Td\n14 TL\n")
	for i, line := range escaped {
		if i > 0 {
			builder.WriteString("T*\n")
		}
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("ET")
	return builder.String()
}
