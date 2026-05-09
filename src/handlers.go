package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultBanner   = "standard"
	defaultTemplate = "index.html"
	defaultWidth    = uint16(80)
	// Input Length Limiting (DoS Protection)
	maxInputLength = 2000
)

type pageData struct {
	Input        string
	Banner       string
	Result       string
	ErrorMessage string
	InputLength  int
}

// renderTemplate renders a template from the templates directory to the response.
func renderTemplate(w http.ResponseWriter, statusCode int, templateFile string, data any) error {
	templatePath := filepath.Join("templates", templateFile)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(body.Bytes())
	return err
}

// renderASCIIArt generates ASCII art with the selected banner using default web options.
func RenderASCIIArt(input, banner string) (string, error) {
	fontData := GetReader(FontType(banner))
	fontMap, err := BuildFontMap(fontData)
	if err != nil {
		return "", err
	}

	options := BannerOptions{Align: "left"}
	return RenderInputStringWeb(input, fontMap, options, defaultWidth), nil
}

func renderTemplateError(w http.ResponseWriter, err error) {
	if errors.Is(err, os.ErrNotExist) {
		http.Error(w, "template not found", http.StatusNotFound)
		return
	}

	http.Error(w, "HTTP status 500 - Internal Server Error", http.StatusInternalServerError)
}

// HomeHandler handles requests for the main page.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	fmt.Println("GET request from: ", r.RemoteAddr)
	if r.Method != http.MethodGet {
		data := pageData{
			Banner:       defaultBanner,
			ErrorMessage: "HTTP status 400 - bad request: GET is required for this route",
		}
		if err := renderTemplate(w, http.StatusBadRequest, defaultTemplate, data); err != nil {
			renderTemplateError(w, err)
		}
		return
	}

	data := pageData{
		Input:        "",
		Banner:       defaultBanner,
		Result:       "",
		ErrorMessage: "",
		InputLength:  0,
	}

	if err := renderTemplate(w, http.StatusOK, defaultTemplate, data); err != nil {
		renderTemplateError(w, err)
	}
}

// ApiAsciiArtHandler handles JSON requests for ASCII art generation.

type asciiAPIRequest struct {
	Text   string `json:"text"`
	Banner string `json:"banner,omitempty"`
	Color  string `json:"color,omitempty"`
}

type asciiAPIResponse struct {
	Input  string `json:"input"`
	Banner string `json:"banner"`
	Color  string `json:"color,omitempty"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func ApiAsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method must be POST",
		})
		return
	}

	defer r.Body.Close()

	var req asciiAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
		return
	}

	input := strings.ReplaceAll(req.Text, "\r\n", "\n")
	banner := req.Banner
	if banner == "" {
		banner = defaultBanner
	}

	if len(input) > maxInputLength {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "input text is too long",
		})
		return
	}

	if !IsValidAscii(input) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "ERROR 400 bad request: input contains non-ASCII characters",
		})
		return
	}

	if !IsValidFont(banner) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid banner",
		})
		return
	}

	result, err := RenderASCIIArt(input, banner)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Error 500: internal server error",
		})
		return
	}

	// Apply color if provided
	if req.Color != "" {
		colorVal, err := ParseHexColor(req.Color)
		if err == nil {
			result = Colorize(result, colorVal)
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"result": result,
		"color":  req.Color,
	})
}
