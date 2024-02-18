package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"linebot/applicationerror"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type applicationLog struct {
	Id        string
	MsgFormat string
}

func newApplicationLog(id, msgFormat string) *applicationLog {
	return &applicationLog{
		Id:        id,
		MsgFormat: msgFormat,
	}
}

var (
	// WebHook署名検証開始
	LBIF010001 = newApplicationLog("LBIF010001", "WebHookの署名を検証します")
	// WebHook署名検証成功
	LBIF010002 = newApplicationLog("LBIF010002", "WebHookの署名の検証に成功しました")
	// 不正なユーザからのリクエスト受信
	LBIF020001 = newApplicationLog("LBIF020001", "不正なユーザからのリクエストです。userId:%s")
	// 不明なユーザからのリクエスト受信
	LBIF020002 = newApplicationLog("LBIF020002", "不明なユーザからのリクエストです。userId:%s")
	// 鍵サーバへの接続{path}
	LBIF040001 = newApplicationLog("LBIF040001", "鍵サーバに接続します、%s")
	// 鍵サーバからレスポンス受信(response)
	LBIF040002 = newApplicationLog("LBIF040002", "鍵サーバからレスポンスを受信しました。%s")
	// Lineへの返信メッセージ応答(ReplyToken,Message)
	LBIF050001 = newApplicationLog("LBIF050001", "メッセージを応答します。ReplyToken:%s,Message:%s")
	// WebHook署名検証エラー
	LBWR010001 = newApplicationLog("LBWR010001", "WebHookの署名の検証中にエラーが発生しました。")
	// WebHook署名検証失敗
	LBWR010002 = newApplicationLog("LBWR010002", "WebHookの署名の検証に失敗しました。")
	// メッセージ応答時にエラーが発生しました。replyToken:%replyToken,message:%message
	LBWR050001 = newApplicationLog("LBWR050001", "メッセージ応答時にエラーが発生しました。replyToken:%s,message:%s")
	// DBでのレコード検索時に指定テーブルにおいて、指定キーのレコードが見つからなかった
	LBER030001 = newApplicationLog("LBER030001", "指定されたレコードが見つかりません。")
	// DBへのレコード保管時にエラー
	LBER030002 = newApplicationLog("LBER030002", "レコード記録時にエラーが発生しました。")
	// Requestログ失敗
	LBER010001 = newApplicationLog("LBER010001", "Requestログに失敗しました")
	// DB接続失敗
	LBFT030001 = newApplicationLog("LBFT030001", "DBとの接続に失敗しました。")
	// 鍵サーバ接続失敗
	LBFT040001 = newApplicationLog("LBFT040001", "鍵サーバとの接続に失敗しました。")
	// 鍵サーバレスポンス読み込み失敗
	LBFT040002 = newApplicationLog("LBFT040002", "鍵サーバのレスポンス読み込みに失敗しました。")
	// 鍵サーバレスポンス形式不正{response}
	LBFT040003 = newApplicationLog("LBFT040003", "鍵サーバのレスポンス形式が不正です。%s")
	// LINE Botの初期化失敗
	LBFT040004 = newApplicationLog("LBFT040004", "LineBotの初期化に失敗しました。")
	// サーバ起動ログ
	LBIF900001 = newApplicationLog("LBIF900001", "Server Initialize Completed : app port=%s, metrics port=%s")
	// 暗号化処理に失敗
	LBFT900001 = newApplicationLog("LBFT900001", "暗号化処理に失敗しました。")
	// 環境変数の読み込みに失敗(環境変数名)
	LBFT900002 = newApplicationLog("LBFT900002", "環境変数の読み込みに失敗しました。%s")
	// システム障害が発生。
	LBFT909999 = newApplicationLog("LBFT909999", "システム障害が発生しました。")
)

func (v applicationLog) GetId() string {
	return v.Id
}

func (v applicationLog) GetMsgFormat() string {
	return v.MsgFormat
}

func Request(req *http.Request) {
	buf, err := io.ReadAll(req.Body)
	if err != nil {

		ErrorWithStackTrace(err, nil, LBER010001)
	}

	j, err := json.Marshal(req.Header)
	if err != nil {
		ErrorWithStackTrace(err, nil, LBER010001)
	}

	log.Info().Str("type", "RequestLogger").RawJSON("header", j).Msg(string(buf))

	// 後続でRequestの内容を読むために詰め直し
	req.Body = io.NopCloser(bytes.NewBuffer(buf))
}

func Debug(message string) {
	if os.Getenv("LOG_LEVEL") == "debug" {
		log.Debug().Str("type", "DebugLogger").Msg(message)
	}
}

func Info(l *applicationLog, values ...interface{}) {
	log.Info().Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func Warn(l *applicationLog, values ...interface{}) {
	log.Warn().Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func WarnWithStackTrace(err error, aplErr *applicationerror.ApplicationError, l *applicationLog, values ...interface{}) {
	if aplErr != nil {
		log.Warn().Err(err).Str("type", "ApplicationLogger").Str("errorCode", aplErr.Code).Str("errorMessage", aplErr.Message).Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
	}
	log.Warn().Err(err).Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func Error(l *applicationLog, values ...interface{}) {
	log.Error().Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func ErrorWithStackTrace(err error, aplErr *applicationerror.ApplicationError, l *applicationLog, values ...interface{}) {
	if aplErr != nil {
		log.Error().Err(err).Str("type", "ApplicationLogger").Str("errorCode", aplErr.Code).Str("errorMessage", aplErr.Message).Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
	}
	log.Error().Err(err).Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func Fatal(l *applicationLog, values ...interface{}) {
	log.WithLevel(zerolog.FatalLevel).Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func FatalWithStackTrace(err error, aplErr *applicationerror.ApplicationError, l *applicationLog, values ...interface{}) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if aplErr != nil {
		log.WithLevel(zerolog.FatalLevel).Err(err).Str("type", "ApplicationLogger").Str("errorCode", aplErr.Code).Str("errorMessage", aplErr.Message).Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
	}
	log.WithLevel(zerolog.FatalLevel).Err(err).Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}
