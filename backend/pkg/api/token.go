package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wakka-2/Namless/backend/pkg/types"
)

const (
	tokenURL = "https://studio-api.nmkr.io/v2/UploadNft/651195c1-4663-451d-90dd-bb3d8cfb8f43?uploadsource=Boom"
)

// CreateToken creates a new token.
//
//nolint:funlen
func (r *RESTAPI) CreateToken(writer http.ResponseWriter, req *http.Request) {
	input := types.TokenInput{}

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	payload := map[string]any{
		"tokenname":   input.Tokenname,
		"displayname": input.Displayname,
		"description": input.Description,
		"previewImageNft": map[string]string{
			"mimetype":       "string",
			"fileFromIPFS":   input.FileFromIPFS,
			"fileFromBase64": input.FileFromBase64,
		},
		"metadataPlaceholder": []map[string]any{
			{
				"name":  input.MetadataPlaceholderName,
				"value": input.MetadataPlaceholderValue,
			},
		},
	}

	asJSON, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	request, err := http.NewRequestWithContext(req.Context(), http.MethodPost, tokenURL, bytes.NewBuffer(asJSON))
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	request.Header = http.Header{
		"Accept":       {"text/plain"},
		"Content-Type": {"application/json"},
	}
	request.Header.Add("Authorization", input.Bearer)

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = write(writer, bytes, http.StatusOK)
	if err != nil {
		log.Default().Printf("could not write: %s", err)
	}

	// TODO: remove
	fmt.Println("## bytes response:", string(bytes))
}

// CreateToken2 creates a new token.
//
//nolint:funlen
func (r *RESTAPI) CreateToken2(writer http.ResponseWriter, req *http.Request) {
	// TODO: remove
	fmt.Println("##222 ")

	key := req.PathValue("name")
	if key == "" {
		r.handleError(writer, "missing key", http.StatusBadRequest)
		return
	}

	// TODO: remove
	fmt.Println("##")
	fmt.Println("## CreateToken2 #### :", key)

	url := fmt.Sprintf("https://studio-api.nmkr.io/v2/MintAndSendSpecific/651195c1-4663-451d-90dd-bb3d8cfb8f43/%s/0.1/addr1qxz79ckr8ewtupsa3reut65t7l80xalv3d8vp38855jkq34hdj2er5qm75szuxaw2aw28qy9wyaywheaxghvn9hgq60qypfa2x",
		key)

	request, err := http.NewRequestWithContext(req.Context(), http.MethodGet, url, nil)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	request.Header = http.Header{
		"Accept": {"text/plain"},
	}
	request.Header.Add("Authorization", "59b826e660db4fbf941aabe3b5d5b84b")

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		r.handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = write(writer, bytes, http.StatusOK)
	if err != nil {
		log.Default().Printf("could not write: %s", err)
	}

	// TODO: remove
	fmt.Println("##")
	fmt.Println("##")
	fmt.Println("##")
	fmt.Println("## bytes response 222:", string(bytes))
}
