package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "google.golang.org/genai"
)

func main() {
    // تحميل API key من environment variable
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        log.Fatal("GEMINI_API_KEY environment variable is required")
    }
    
    ctx := context.Background()
    
    // إنشاء client بالطريقة الصحيحة
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: apiKey,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 1. توليد نص (مثل الكود الأصلي)
    fmt.Println("=== توليد نص ===")
    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash",
        genai.Text("Explain how AI works in a few words"),
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // طباعة النتيجة
    for _, part := range result.Candidates[0].Content.Parts {
        if part.Text != "" {
            fmt.Println(part.Text)
        }
    }
    
    // 2. توليد صورة (مثل كود nano banana)
    fmt.Println("\n=== توليد صورة ===")
    imageResult, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash-image",
        genai.Text("Create a picture of a nano banana dish in a fancy restaurant with a Gemini theme"),
    )
    if err != nil {
        log.Printf("Image generation error: %v", err)
    } else {
        // حفظ الصورة
        for _, part := range imageResult.Candidates[0].Content.Parts {
            if part.InlineData != nil {
                imageBytes := part.InlineData.Data
                outputFilename := "gemini_generated_image.png"
                if err := os.WriteFile(outputFilename, imageBytes, 0644); err != nil {
                    log.Printf("Failed to save image: %v", err)
                } else {
                    fmt.Printf("✅ Image saved as %s\n", outputFilename)
                }
            }
        }
    }
}