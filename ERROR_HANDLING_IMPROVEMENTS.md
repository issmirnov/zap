# Error Handling Improvements for Zap

This document outlines the comprehensive improvements made to error handling in the Zap project, addressing GitHub issue #3: "Refactor error handling to send helpful 500's".

## Overview

The goal was to improve error handling throughout the application to provide users with more helpful error messages instead of generic failures. Previously, many errors were silently ignored or provided minimal information, which could confuse users.

## Key Improvements

### 1. Enhanced HTTP Error Responses

**File: `cmd/zap/structs.go`**
- Improved `CtxWrapper.ServeHTTP()` to provide specific error messages for different HTTP status codes
- Added helpful prefixes for 500, 404, and 400 errors
- Better error message formatting for user clarity

**Before:**
```go
http.Error(w, fmt.Sprintf("HTTP %d: %q", status, err), status)
```

**After:**
```go
case http.StatusInternalServerError:
    errorMsg := fmt.Sprintf("Internal Server Error: %s", err.Error())
    http.Error(w, errorMsg, status)
case http.StatusNotFound:
    errorMsg := fmt.Sprintf("Shortcut not found: %s", err.Error())
    http.Error(w, errorMsg, status)
```

### 2. Improved Request Handler Error Handling

**File: `cmd/zap/web.go`**
- Enhanced `IndexHandler` with comprehensive validation checks
- Added nil configuration validation
- Better error messages for configuration structure issues
- Improved schema validation error handling
- Added path generation validation

**New validations:**
```go
// Check if context and configuration are valid
if ctx == nil || ctx.Config == nil {
    return http.StatusInternalServerError, fmt.Errorf("server configuration is invalid or not loaded")
}

// Validate configuration structure
if conf == nil {
    return http.StatusInternalServerError, fmt.Errorf("invalid configuration structure for host '%s'", host)
}

// Validate path generation
if path.Len() == 0 {
    return http.StatusInternalServerError, fmt.Errorf("failed to generate redirect path for host '%s'", host)
}
```

### 3. Enhanced Configuration Validation

**File: `cmd/zap/config.go`**
- Improved `ValidateConfig()` function with more descriptive error messages
- Better type checking with clear error descriptions
- Enhanced error context for nested validation failures
- Added nil configuration validation

**Improved error messages:**
```go
// Before
errors = multierror.Append(errors, fmt.Errorf("expected bool value for %T, got: %v", k, v.Data()))

// After
errors = multierror.Append(errors, fmt.Errorf("expected boolean value for 'ssl_off' key, got: %T (%v)", v.Data(), v.Data()))
```

### 4. Better Configuration Parsing

**File: `cmd/zap/config.go`**
- Enhanced `ParseYaml()` and `parseYamlString()` functions
- Added input validation (empty strings, nil inputs)
- Better error wrapping with context
- Improved error messages for file operations

**New validations:**
```go
if fname == "" {
    return nil, fmt.Errorf("no configuration file specified")
}

if len(data) == 0 {
    return nil, fmt.Errorf("configuration file '%s' is empty", fname)
}
```

### 5. Improved Path Expansion Error Handling

**File: `cmd/zap/text.go`**
- Enhanced `ExpandPath()` function to return errors
- Better error context for path expansion failures
- Improved validation of configuration structure during expansion
- Replaced panic with proper error handling

**Before:**
```go
func ExpandPath(c *gabs.Container, token *list.Element, res *bytes.Buffer) {
    expandPath(c, token, res, true)
}
```

**After:**
```go
func ExpandPath(c *gabs.Container, token *list.Element, res *bytes.Buffer) error {
    return expandPath(c, token, res, true)
}
```

### 6. Enhanced Main Application Error Handling

**File: `cmd/main.go`**
- Better error messages for configuration failures
- Improved validation feedback
- Enhanced server startup information
- Better error handling for file watcher creation
- Graceful handling of hosts file update failures

**New startup information:**
```go
fmt.Printf("Configuration file: %s\n", *configName)
fmt.Printf("Health check: http://%s/healthz\n", serverAddr)
fmt.Printf("Configuration view: http://%s/varz\n", serverAddr)
```

### 7. Improved File System Operations

**File: `cmd/zap/config.go`**
- Enhanced `UpdateHosts()` function to return errors
- Better error handling for file read/write operations
- Improved error messages for permission issues
- Better logging for configuration reload operations

### 8. Enhanced Testing

**File: `cmd/zap/web_test.go`**
- Added new test cases for error handling scenarios
- Updated existing tests to reflect improved error messages
- Better test coverage for error conditions

## Error Message Examples

### Before (Generic)
```
HTTP 500: "failed to write response"
```

### After (Helpful)
```
Internal Server Error: failed to write response
```

### Before (Minimal)
```
shortcut 'example.com' not found in config
```

### After (Clear)
```
Shortcut not found: shortcut 'example.com' not found in config
```

### Before (Unclear)
```
expected bool value for string, got: not_bool
```

### After (Descriptive)
```
expected boolean value for 'ssl_off' key, got: string (not_bool)
```

## Benefits

1. **Better User Experience**: Users now receive clear, actionable error messages
2. **Easier Debugging**: Developers can quickly identify configuration issues
3. **Improved Reliability**: Better error handling prevents silent failures
4. **Enhanced Monitoring**: Clear error messages make it easier to monitor application health
5. **Better Documentation**: Error messages serve as implicit documentation of expected behavior

## Backward Compatibility

All improvements maintain backward compatibility:
- No breaking changes to public APIs
- Existing functionality preserved
- Enhanced error handling is additive

## Testing

All improvements include comprehensive testing:
- Unit tests for error conditions
- Integration tests for error handling flows
- Updated existing tests to reflect new error message formats

## Future Enhancements

Potential areas for further improvement:
1. Structured error logging (JSON format)
2. Error code system for programmatic error handling
3. User-friendly error pages for web interface
4. Error reporting and analytics
5. Configuration validation at startup with detailed feedback

## Conclusion

These improvements transform Zap from having minimal error handling to providing comprehensive, user-friendly error messages. Users can now quickly understand what went wrong and how to fix it, while developers have better tools for debugging and monitoring the application.
