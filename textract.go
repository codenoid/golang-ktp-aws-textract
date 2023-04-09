package main

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/textract"
)

func AnalyzeImage(client *textract.Textract, imageBytes []byte) (map[string]string, error) {
	input := &textract.AnalyzeDocumentInput{
		Document: &textract.Document{
			Bytes: imageBytes,
		},
		FeatureTypes: []*string{
			aws.String("FORMS"),
		},
	}
	output, err := client.AnalyzeDocument(input)
	if err != nil {
		return nil, err
	}

	blockMap := map[string]textract.Block{}
	for _, block := range output.Blocks {
		blockMap[*block.Id] = *block
	}

	keyValues := make([]map[string]string, 0)
	for _, block := range output.Blocks {
		if *block.BlockType == "KEY_VALUE_SET" {
			if len(block.EntityTypes) > 0 {
				if *block.EntityTypes[0] == "KEY" {
					keyValues = append(keyValues, map[string]string{
						GetChild(*block, blockMap): GetValue(*block, blockMap),
					})
				}
			}
		}
	}

	kvval := NormalizeKTPKey(keyValues)
	return kvval, nil
}

func GetChild(block textract.Block, blockMap map[string]textract.Block) string {
	for _, id := range block.Relationships {
		if *id.Type == "CHILD" {
			childText := ""
			for _, childID := range id.Ids {
				if val, ok := blockMap[*childID]; ok {
					if val.Text != nil {
						childText += *val.Text + " "
					}
				}
			}
			return StripNonAlphaNumSpace(childText)
		}
	}
	return ""
}

func GetValue(block textract.Block, blockMap map[string]textract.Block) string {
	for _, id := range block.Relationships {
		if *id.Type == "VALUE" {
			childText := ""
			for _, childID := range id.Ids {
				if val, ok := blockMap[*childID]; ok {
					if val.Text != nil {
						childText += *val.Text + " "
					}
					childText += GetChild(val, blockMap)
				}
			}
			return StripNonAlphaNumSpace(childText)
		}
	}
	return ""
}

func StripNonNumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if '0' <= b && b <= '9' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func StripNonAlphaNumSpace(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' || b == '/' ||
			b == '-' || b == ',' ||
			b == '.' {
			result.WriteByte(b)
		}
	}
	return strings.TrimSpace(result.String())
}
