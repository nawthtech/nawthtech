package ai

// VideoServiceInterface واجهة لفصل التبعيات
type VideoServiceInterface interface {
    GenerateVideo(prompt string, options VideoOptions) (*VideoResponse, error)
    // ... other methods
}

type VideoOptions struct {
    Duration int
    Quality  string
    // ... other fields
}

type VideoResponse struct {
    URL      string
    Duration int
    // ... other fields
}