package agent

// That fallback hash is used only to avoid time-based attacks
// We want a valid hash to be here since we actually want the bcrypt compare to run
// even if afterwards we will discard its result
const bcryptFallbackHash = "$2y$10$ulOzTkNjLpMJI2HS7rq5/eEIX6qLa9JtcKQ9uk8PjYq68ZUsMN5di "
const invalidUserId = "invalid user id"
