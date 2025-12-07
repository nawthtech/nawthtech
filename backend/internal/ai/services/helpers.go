package services

import (
    "context"
)

func extractUserIDFromContext(ctx context.Context) string {
    if userID, ok := ctx.Value("user_id").(string); ok {
        return userID
    }
    return ""
}

func extractUserTierFromContext(ctx context.Context) string {
    if userTier, ok := ctx.Value("user_tier").(string); ok {
        return userTier
    }
    return "free"
}