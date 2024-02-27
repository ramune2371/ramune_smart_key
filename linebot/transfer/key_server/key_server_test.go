package key_server_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"linebot/applicationerror"
	"linebot/entity"
	"linebot/props"
	"linebot/testutil"
	"linebot/transfer/key_server"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockedClient struct{}
type mockedReadCloser struct{}

func (mockedReadCloser) Read(p []byte) (int, error) {
	return 0, errors.New("test")
}

func (mockedReadCloser) Close() error {
	return nil
}

func (mockedClient) Do(c *http.Request) (*http.Response, error) {
	return &http.Response{Body: mockedReadCloser{}}, nil
}

type URLCheckClient struct {
	test       *testing.T
	expectPath string
}

func (ucc URLCheckClient) Do(c *http.Request) (*http.Response, error) {
	if c.URL.Path != ucc.expectPath {
		ucc.test.Errorf(testutil.STRING_TEST_MSG_FMT, "", ucc.expectPath, c.URL.Path)
	}
	return &http.Response{Body: io.NopCloser(bytes.NewReader([]byte(`{"test":"test"}`)))}, nil
}

// テスト用のHTTPサーバーを起動し、リクエストに応答する関数
func startTestServer(statusCode int, responseBody string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))
	props.KeyServerURL = server.URL // テストサーバーのURLを設定
	return server
}

func createResponse(keyStatus, operationStatus string) string {
	return fmt.Sprintf(`{"keyStatus":"%s", "operationStatus":"%s"}`, keyStatus, operationStatus)
}

func TestKeyServerTransferImpl_Request(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		// テスト用のHTTPサーバーを起動
		oldKeyServerURL := props.KeyServerURL
		defer func() { props.KeyServerURL = oldKeyServerURL }()
		testServer := startTestServer(http.StatusOK, createResponse("True", "another"))
		defer testServer.Close()

		// テスト用のHTTPクライアントを作成
		testClient := testServer.Client()

		// テスト用のKeyServerTransferImplを作成
		keyServerTransfer := key_server.NewKeyServerTransferImpl(testClient)

		// テストを実行
		response, err := keyServerTransfer.Request("")

		// レスポンスを検証
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expect := entity.KeyServerResponse{KeyStatus: entity.KeyStatusOpen, OperationStatus: entity.OperationAnother}
		if response != expect {
			t.Errorf("Unexpected response. Expected: %v, got: %v", expect, response)
		}
	})

	t.Run("異常系(status 200以外)", func(t *testing.T) {
		// エラーレスポンスのステータスコードとメッセージ
		testErrorCode := http.StatusInternalServerError
		testErrorMessage := "Internal Server Error"
		// テスト用のHTTPサーバーを起動
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, testErrorMessage, testErrorCode)
		}))
		defer testServer.Close()

		// テスト用のHTTPクライアントを作成
		testClient := testServer.Client()

		// テスト用のKeyServerTransferImplを作成
		keyServerTransfer := key_server.NewKeyServerTransferImpl(testClient)

		// テストを実行
		response, err := keyServerTransfer.Request("")

		// エラーを検証
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, applicationerror.ConnectionError) {
			t.Errorf("Expected connection error, got: %v", err)
		}
		if response != (entity.KeyServerResponse{}) {
			t.Errorf("Expected empty response, got: %v", response)
		}
	})

	t.Run("異常系(response読み込みエラー)", func(t *testing.T) {

		// テスト用のKeyServerTransferImplを作成
		keyServerTransfer := key_server.NewKeyServerTransferImpl(mockedClient{})

		// テストを実行
		response, err := keyServerTransfer.Request("")

		// エラーを検証
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, applicationerror.ResponseParseError) {
			t.Errorf("Expected ResponseParse error, got: %v", err)
		}
		if response != (entity.KeyServerResponse{}) {
			t.Errorf("Expected empty response, got: %v", response)
		}
	})

	t.Run("異常系(レスポンスが非json)", func(t *testing.T) {
		// テスト用のHTTPサーバーを起動
		oldKeyServerURL := props.KeyServerURL
		defer func() { props.KeyServerURL = oldKeyServerURL }()
		testServer := startTestServer(http.StatusOK, "not json")
		defer testServer.Close()

		// テスト用のHTTPクライアントを作成
		testClient := testServer.Client()

		// テスト用のKeyServerTransferImplを作成
		keyServerTransfer := key_server.NewKeyServerTransferImpl(testClient)

		// テストを実行
		response, err := keyServerTransfer.Request("")

		// エラーを検証
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !errors.Is(err, applicationerror.ResponseParseError) {
			t.Errorf("Expected connection error, got: %v", err)
		}
		if response != (entity.KeyServerResponse{}) {
			t.Errorf("Expected empty response, got: %v", response)
		}
	})
}

func TestCheckKey(t *testing.T) {
	client := URLCheckClient{test: t, expectPath: "/check"}
	oldKeyServerURL := props.KeyServerURL
	defer func() { props.KeyServerURL = oldKeyServerURL }()
	props.KeyServerURL = "http://localhost:8282/"
	kst := key_server.NewKeyServerTransferImpl(&client)
	kst.CheckKey()
}

func TestOpenKey(t *testing.T) {
	client := URLCheckClient{test: t, expectPath: "/open"}
	oldKeyServerURL := props.KeyServerURL
	defer func() { props.KeyServerURL = oldKeyServerURL }()
	props.KeyServerURL = "http://localhost:8282/"
	kst := key_server.NewKeyServerTransferImpl(client)
	kst.OpenKey()
}

func TestCloseKey(t *testing.T) {
	client := URLCheckClient{test: t, expectPath: "/close"}
	oldKeyServerURL := props.KeyServerURL
	defer func() { props.KeyServerURL = oldKeyServerURL }()
	props.KeyServerURL = "http://localhost:8282/"
	kst := key_server.NewKeyServerTransferImpl(client)
	kst.CloseKey()
}
