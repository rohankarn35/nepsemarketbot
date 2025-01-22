package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/rohankarn35/nepsemarketbot/models"
)

func FetchFPOData(url string) (*models.FpoRoot, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var wg sync.WaitGroup
	var body []byte
	var readErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		body, readErr = ioutil.ReadAll(resp.Body)
	}()

	wg.Wait()
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response body: %v", readErr)
	}

	var root models.FpoRoot
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &root, nil
}
