package middleware

/*
// User represents a user with roles.
type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}


func generateJWT(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"roles": roles,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func retrieveJWTToken(r *http.Request) (string, error) {
	// Try to retrieve the token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString != "" {
		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			return tokenString[7:], nil
		}
		return tokenString, nil
	}

	// Try to retrieve the token from a cookie named "jwt_token"
	cookie, err := r.Cookie("jwt_token")
	if err == nil && cookie != nil {
		return cookie.Value, nil
	}

	return "", fmt.Errorf("Token not found")
}

func authorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the JWT token from header or cookie
		tokenString, err := retrieveJWTToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store claims in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if the user has the required role to access this resource
	requiredRole := "admin" // Change to the desired role
	roles, ok := claims["roles"].([]interface{})
	if !ok || !hasRole(roles, requiredRole) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Fetch and return the user's profile data
	userID, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userProfile, err := getUserProfileData(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	// Send the user profile as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func hasRole(roles []interface{}, roleToCheck string) bool {
	for _, role := range roles {
		if role == roleToCheck {
			return true
		}
	}
	return false
}

func getUserProfileData(userID string) (map[string]interface{}, error) {
	// Fetch user profile data from your data source
	// Replace this with your actual data retrieval logic

	userProfile := map[string]interface{}{
		"userID": userID,
		"name":   "John Doe",
		"email":  "johndoe@example.com",
	}

	return userProfile, nil
}
*/
