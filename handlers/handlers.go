package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type HTMLRequest struct {
	HTML string `json:"html"`
}

func DownloadPDF() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Downloading PDF...")

		// Parse request body
		var req HTMLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error parsing request: %v", err)
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		// Create Chrome instance
		ctx, cancel := chromedp.NewContext(
			context.Background(),
			chromedp.WithLogf(log.Printf),
		)
		defer cancel()

		// Add timeout for server response
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// Create a buffer to store the PDF
		var buf []byte

		// Navigate and generate PDF
		if err := chromedp.Run(ctx,
			chromedp.Navigate("about:blank"),
			chromedp.ActionFunc(func(ctx context.Context) error {
				frameTree, err := page.GetFrameTree().Do(ctx)
				if err != nil {
					return err
				}
				return page.SetDocumentContent(frameTree.Frame.ID, req.HTML).Do(ctx)
			}),
			//Recreate the loaded content with all styles in an A4 format  and generate the pdf
			chromedp.ActionFunc(func(ctx context.Context) error {
				var err error
				buf, _, err = page.PrintToPDF().
					WithPrintBackground(true).
					WithPreferCSSPageSize(true).
					WithPaperWidth(8.27).
					WithPaperHeight(11.7).
					WithMarginTop(0).
					WithMarginBottom(0).
					WithMarginLeft(0).
					WithMarginRight(0).
					WithScale(1.0).
					Do(ctx)
				return err
			}),
		); err != nil {
			log.Printf("Error generating PDF: %v", err)
			http.Error(w, fmt.Sprintf("Failed to generate PDF: %v", err), http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")

		// Send the PDF
		w.Write(buf)
	}
}
