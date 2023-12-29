package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

var (
	// WebHook署名検証開始
	LBIF010001 = applicationLog{Id: "LBIF010001", MsgFormat: "WebHookの署名を検証します"}
	// WebHook署名検証成功
	LBIF010002 = applicationLog{Id: "LBIF010002", MsgFormat: "WebHookの署名の検証に成功しました"}
	// 不正なユーザからのリクエスト受信
	LBIF020001 = applicationLog{Id: "LBIF020001", MsgFormat: "不正なユーザからのリクエストです。userId:%s"}
	// 鍵サーバへの接続{path}
	LBIF040001 = applicationLog{Id: "LBIF040001", MsgFormat: "鍵サーバに接続します、%s"}
	// 鍵サーバからレスポンス受信(response)
	LBIF040002 = applicationLog{Id: "LBIF040002", MsgFormat: "鍵サーバからレスポンスを受信しました。%s"}
	// Lineへの返信メッセージ応答(ReplyToken,Message)
	LBIF050001 = applicationLog{Id: "LBIF050001", MsgFormat: "メッセージを応答します。ReplyToken:%s,Message:%s"}
	// WebHook署名検証エラー
	LBWR010001 = applicationLog{Id: "LBWR010001", MsgFormat: "WebHookの署名の検証中にエラーが発生しました。"}
	// WebHook署名検証失敗
	LBWR010002 = applicationLog{Id: "LBWR010002", MsgFormat: "WebHookの署名の検証に失敗しました。"}
	// メッセージ応答時にエラーが発生しました。replyToken:%replyToken,message:%message
	LBWR050001 = applicationLog{Id: "LBWR050001", MsgFormat: "メッセージ応答時にエラーが発生しました。replyToken:%s,message:%s"}
	// DBでのレコード検索時に指定テーブルにおいて、指定キーのレコードが見つからなかった
	LBER030001 = applicationLog{Id: "LBER030001", MsgFormat: "指定されたレコードが見つかりません。"}
	// DBへのレコード保管時にエラー
	LBER030002 = applicationLog{Id: "LBER030002", MsgFormat: "レコード記録時にエラーが発生しました。"}
	// Requestログ失敗
	LBER010001 = applicationLog{Id: "LBER010001", MsgFormat: "Requestログに失敗しました"}
	// DB接続失敗
	LBFT030001 = applicationLog{Id: "LBFT030001", MsgFormat: "DBとの接続に失敗しました。"}
	// 鍵サーバ接続失敗
	LBFT040001 = applicationLog{Id: "LBFT040001", MsgFormat: "鍵サーバとの接続に失敗しました。"}
	// 鍵サーバレスポンス読み込み失敗
	LBFT040002 = applicationLog{Id: "LBFT040002", MsgFormat: "鍵サーバのレスポンス読み込みに失敗しました。"}
	// 鍵サーバレスポンス形式不正{response}
	LBFT040003 = applicationLog{Id: "LBFT040003", MsgFormat: "鍵サーバのレスポンス形式が不正です。%s"}
	// LINE Botの初期化失敗
	LBFT040004 = applicationLog{Id: "LBFT040004", MsgFormat: "LineBotの初期化に失敗しました。"}
	// サーバ起動ログ
	LBIF900001 = applicationLog{Id: "LBIF900001", MsgFormat: "Server Initialize Completed : app port=%s, metrics port=%s"}
	// 暗号化処理に失敗
	LBFT900001 = applicationLog{Id: "LBFT900001", MsgFormat: "暗号化処理に失敗しました。"}
	// システム障害が発生。
	LBFT909999 = applicationLog{Id: "LBFT909999", MsgFormat: "システム障害が発生しました。"}
)

func (v *applicationLog) GetId() string {
	return v.Id
}

func (v *applicationLog) GetMsgFormat() string {
	return v.MsgFormat
}

func Request(req *http.Request) {
	buf, err := io.ReadAll(req.Body)
	if err != nil {

		ErrorWithStackTrace(err, &LBER010001)
	}

	j, err := json.Marshal(req.Header)
	if err != nil {
		ErrorWithStackTrace(err, &LBER010001)
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

func WarnWithStackTrace(err error, l *applicationLog, values ...interface{}) {
	log.Warn().Err(err).Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func Error(l *applicationLog, values ...interface{}) {
	log.Error().Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func ErrorWithStackTrace(err error, l *applicationLog, values ...interface{}) {
	log.Error().Err(err).Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func Fatal(l *applicationLog, values ...interface{}) {
	log.WithLevel(zerolog.FatalLevel).Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}

func FatalWithStackTrace(err error, l *applicationLog, values ...interface{}) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.WithLevel(zerolog.FatalLevel).Err(err).Str("type", "ApplicationLogger").Str("id", l.Id).Msg(fmt.Sprintf(l.GetMsgFormat(), values...))
}
