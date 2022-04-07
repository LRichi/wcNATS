package nats

import (
	"strings"

	"github.com/nats-io/nats.go"
)

type data struct {
	m string
}

func (e data) Error() string {
	return e.m
}

type ConnectionClosed struct {
	data //  "nats: connection closed"
}

type ConnectionDraining struct {
	data //  "nats: connection draining"
}

type DrainTimeout struct {
	data //  "nats: draining connection timed out"
}

type ConnectionReconnecting struct {
	data //  "nats: connection reconnecting"
}

type SecureConnRequired struct {
	data //  "nats: secure connection required"
}

type SecureConnWanted struct {
	data //  "nats: secure connection not available"
}

type BadSubscription struct {
	data //  "nats: invalid subscription"
}

type TypeSubscription struct {
	data //  "nats: invalid subscription type"
}

type BadSubject struct {
	data //  "nats: invalid subject"
}

type BadQueueName struct {
	data //  "nats: invalid queue name"
}

type SlowConsumer struct {
	data //  "nats: slow consumer, messages dropped"
}

type Timeout struct {
	data //  "nats: timeout"
}

type BadTimeout struct {
	data //  "nats: timeout invalid"
}

type Authorization struct {
	data //  "nats: authorization violation"
}

type AuthExpired struct {
	data //  "nats: authentication expired"
}

type AuthRevoked struct {
	data //  "nats: authentication revoked"
}

type AccountAuthExpired struct {
	data //  "nats: account authentication expired"
}

type NoServers struct {
	data //  "nats: no servers available for connection"
}

type JsonParse struct {
	data //  "nats: connect message, json parse or"
}

type ChanArg struct {
	data //  "nats: argument needs to be a channel type"
}

type MaxPayload struct {
	data //  "nats: maximum payload exceeded"
}

type MaxMessages struct {
	data //  "nats: maximum messages delivered"
}

type SyncSubRequired struct {
	data //  "nats: illegal call on an async subscription"
}

type MultipleTLSConfigs struct {
	data //  "nats: multiple tls.Configs not allowed"
}

type NoInfoReceived struct {
	data //  "nats: protocol exception, INFO not received"
}

type ReconnectBufExceeded struct {
	data //  "nats: outbound buffer limit exceeded"
}

type InvalidConnection struct {
	data //  "nats: invalid connection"
}

type InvalidMsg struct {
	data //  "nats: invalid message or message nil"
}

type InvalidArg struct {
	data //  "nats: invalid argument"
}

type InvalidContext struct {
	data //  "nats: invalid context"
}

type NoDeadlineContext struct {
	data //  "nats: context requires a deadline"
}

type NoEchoNotSupported struct {
	data //  "nats: no echo option not supported by this server"
}

type ClientIDNotSupported struct {
	data //  "nats: client ID not supported by this server"
}

type UserButNoSigCB struct {
	data //  "nats: user callback defined without a signature handler"
}

type NkeyButNoSigCB struct {
	data //  "nats: nkey defined without a signature handler"
}

type NoUserCB struct {
	data //  "nats: user callback not defined"
}

type NkeyAndUser struct {
	data //  "nats: user callback and nkey defined"
}

type NkeysNotSupported struct {
	data //  "nats: nkeys not supported by the server"
}

type StaleConnection struct {
	data //  "nats: " + STALE_CONNECTION
}

type TokenAlreadySet struct {
	data //  "nats: token and token handler both set"
}

type MsgNotBound struct {
	data //  "nats: message is not bound to subscription/connection"
}

type MsgNoReply struct {
	data //  "nats: message does not have a reply"
}

type ClientIPNotSupported struct {
	data //  "nats: client IP not supported by this server"
}

type Disconnected struct {
	data //  "nats: server is disconnected"
}

type HeadersNotSupported struct {
	data //  "nats: headers not supported by this server"
}

type BadHeaderMsg struct {
	data //  "nats: message could not decode headers"
}

type NoResponders struct {
	data //  "nats: no responders available for request"
}

type NoContextOrTimeout struct {
	data //  "nats: no context or timeout given"
}

type PullModeNotAllowed struct {
	data //  "nats: pull based not supported"
}

type JetStreamNotEnabled struct {
	data //  "nats: jetstream not enabled"
}

type JetStreamBadPre struct {
	data //  "nats: jetstream api prefix not valid"
}

type NoStreamResponse struct {
	data //  "nats: no response from stream"
}

type NotJSMessage struct {
	data //  "nats: not a jetstream message"
}

type InvalidStreamName struct {
	data //  "nats: invalid stream name"
}

type InvalidDurableName struct {
	data //  "nats: invalid durable name"
}

type NoMatchingStream struct {
	data //  "nats: no stream matches subject"
}

type SubjectMismatch struct {
	data //  "nats: subject does not match consumer"
}

type ContextAndTimeout struct {
	data //  "nats: context and timeout can not both be set"
}

type InvalidJSAck struct {
	data //  "nats: invalid jetstream publish response"
}

type MultiStreamUnsupported struct {
	data //  "nats: multiple streams are not supported"
}

type StreamNameRequired struct {
	data //  "nats: stream name is required"
}

type StreamNotFound struct {
	data //  "nats: stream not found"
}

type ConsumerNotFound struct {
	data //  "nats: consumer not found"
}

type ConsumerNameRequired struct {
	data //  "nats: consumer name is required"
}

type ConsumerConfigRequired struct {
	data //  "nats: consumer configuration is required"
}

type StreamSnapshotConfigRequired struct {
	data //  "nats: stream snapshot configuration is required"
}

type DeliverSubjectRequired struct {
	data //  "nats: deliver subject is required"
}

type PullSubscribeToPushConsumer struct {
	data //  "nats: cannot pull subscribe to push based consumer"
}

type PullSubscribeRequired struct {
	data //  "nats: must use pull subscribe to bind to pull based consumer"
}

type ConsumerNotActive struct {
	data //  "nats: consumer not active"
}

type MsgNotFound struct {
	data //  "nats: message not found"
}

func convertErr(err error) error {
	if err == nil {
		return err
	}

	switch err.Error() {
	case "nats: connection closed":
		return ConnectionClosed{data: data{m: convertString(err)}}
	case "nats: connection draining":
		return ConnectionDraining{data: data{m: convertString(err)}}
	case "nats: draining connection timed out":
		return DrainTimeout{data: data{m: convertString(err)}}
	case "nats: connection reconnecting":
		return ConnectionReconnecting{data: data{m: convertString(err)}}
	case "nats: secure connection required":
		return SecureConnRequired{data: data{m: convertString(err)}}
	case "nats: secure connection not available":
		return SecureConnWanted{data: data{m: convertString(err)}}
	case "nats: invalid subscription":
		return BadSubscription{data: data{m: convertString(err)}}
	case "nats: invalid subscription type":
		return TypeSubscription{data: data{m: convertString(err)}}
	case "nats: invalid subject":
		return BadSubject{data: data{m: convertString(err)}}
	case "nats: invalid queue name":
		return BadQueueName{data: data{m: convertString(err)}}
	case "nats: slow consumer, messages dropped":
		return SlowConsumer{data: data{m: convertString(err)}}
	case "nats: timeout":
		return Timeout{data: data{m: convertString(err)}}
	case "nats: timeout invalid":
		return BadTimeout{data: data{m: convertString(err)}}
	case "nats: authorization violation":
		return Authorization{data: data{m: convertString(err)}}
	case "nats: authentication expired":
		return AuthExpired{data: data{m: convertString(err)}}
	case "nats: authentication revoked":
		return AuthRevoked{data: data{m: convertString(err)}}
	case "nats: account authentication expired":
		return AccountAuthExpired{data: data{m: convertString(err)}}
	case "nats: no servers available for connection":
		return NoServers{data: data{m: convertString(err)}}
	case "nats: connect message, json parse error":
		return JsonParse{data: data{m: convertString(err)}}
	case "nats: argument needs to be a channel type":
		return ChanArg{data: data{m: convertString(err)}}
	case "nats: maximum payload exceeded":
		return MaxPayload{data: data{m: convertString(err)}}
	case "nats: maximum messages delivered":
		return MaxMessages{data: data{m: convertString(err)}}
	case "nats: illegal call on an async subscription":
		return SyncSubRequired{data: data{m: convertString(err)}}
	case "nats: multiple tls.Configs not allowed":
		return MultipleTLSConfigs{data: data{m: convertString(err)}}
	case "nats: protocol exception, INFO not received":
		return NoInfoReceived{data: data{m: convertString(err)}}
	case "nats: outbound buffer limit exceeded":
		return ReconnectBufExceeded{data: data{m: convertString(err)}}
	case "nats: invalid connection":
		return InvalidConnection{data: data{m: convertString(err)}}
	case "nats: invalid message or message nil":
		return InvalidMsg{data: data{m: convertString(err)}}
	case "nats: invalid argument":
		return InvalidArg{data: data{m: convertString(err)}}
	case "nats: invalid context":
		return InvalidContext{data: data{m: convertString(err)}}
	case "nats: context requires a deadline":
		return NoDeadlineContext{data: data{m: convertString(err)}}
	case "nats: no echo option not supported by this server":
		return NoEchoNotSupported{data: data{m: convertString(err)}}
	case "nats: client ID not supported by this server":
		return ClientIDNotSupported{data: data{m: convertString(err)}}
	case "nats: user callback defined without a signature handler":
		return UserButNoSigCB{data: data{m: convertString(err)}}
	case "nats: nkey defined without a signature handler":
		return NkeyButNoSigCB{data: data{m: convertString(err)}}
	case "nats: user callback not defined":
		return NoUserCB{data: data{m: convertString(err)}}
	case "nats: user callback and nkey defined":
		return NkeyAndUser{data: data{m: convertString(err)}}
	case "nats: nkeys not supported by the server":
		return NkeysNotSupported{data: data{m: convertString(err)}}
	case "nats: " + nats.STALE_CONNECTION:
		return StaleConnection{data: data{m: convertString(err)}}
	case "nats: token and token handler both set":
		return TokenAlreadySet{data: data{m: convertString(err)}}
	case "nats: message is not bound to subscription/connection":
		return MsgNotBound{data: data{m: convertString(err)}}
	case "nats: message does not have a reply":
		return MsgNoReply{data: data{m: convertString(err)}}
	case "nats: client IP not supported by this server":
		return ClientIPNotSupported{data: data{m: convertString(err)}}
	case "nats: server is disconnected":
		return Disconnected{data: data{m: convertString(err)}}
	case "nats: headers not supported by this server":
		return HeadersNotSupported{data: data{m: convertString(err)}}
	case "nats: message could not decode headers":
		return BadHeaderMsg{data: data{m: convertString(err)}}
	case "nats: no responders available for request":
		return NoResponders{data: data{m: convertString(err)}}
	case "nats: no context or timeout given":
		return NoContextOrTimeout{data: data{m: convertString(err)}}
	case "nats: pull based not supported":
		return PullModeNotAllowed{data: data{m: convertString(err)}}
	case "nats: jetstream not enabled":
		return JetStreamNotEnabled{data: data{m: convertString(err)}}
	case "nats: jetstream api prefix not valid":
		return JetStreamBadPre{data: data{m: convertString(err)}}
	case "nats: no response from stream":
		return NoStreamResponse{data: data{m: convertString(err)}}
	case "nats: not a jetstream message":
		return NotJSMessage{data: data{m: convertString(err)}}
	case "nats: invalid stream name":
		return InvalidStreamName{data: data{m: convertString(err)}}
	case "nats: invalid durable name":
		return InvalidDurableName{data: data{m: convertString(err)}}
	case "nats: no stream matches subject":
		return NoMatchingStream{data: data{m: convertString(err)}}
	case "nats: subject does not match consumer":
		return SubjectMismatch{data: data{m: convertString(err)}}
	case "nats: context and timeout can not both be set":
		return ContextAndTimeout{data: data{m: convertString(err)}}
	case "nats: invalid jetstream publish response":
		return InvalidJSAck{data: data{m: convertString(err)}}
	case "nats: multiple streams are not supported":
		return MultiStreamUnsupported{data: data{m: convertString(err)}}
	case "nats: stream name is required":
		return StreamNameRequired{data: data{m: convertString(err)}}
	case "nats: stream not found":
		return StreamNotFound{data: data{m: convertString(err)}}
	case "nats: consumer not found":
		return ConsumerNotFound{data: data{m: convertString(err)}}
	case "nats: consumer name is required":
		return ConsumerNameRequired{data: data{m: convertString(err)}}
	case "nats: consumer configuration is required":
		return ConsumerConfigRequired{data: data{m: convertString(err)}}
	case "nats: stream snapshot configuration is required":
		return StreamSnapshotConfigRequired{data: data{m: convertString(err)}}
	case "nats: deliver subject is required":
		return DeliverSubjectRequired{data: data{m: convertString(err)}}
	case "nats: cannot pull subscribe to push based consumer":
		return PullSubscribeToPushConsumer{data: data{m: convertString(err)}}
	case "nats: must use pull subscribe to bind to pull based consumer":
		return PullSubscribeRequired{data: data{m: convertString(err)}}
	case "nats: consumer not active":
		return ConsumerNotActive{data: data{m: convertString(err)}}
	case "nats: message not found":
		return MsgNotFound{data: data{m: convertString(err)}}
	default:
		return err
	}
}

func convertString(err error) string {
	return strings.TrimPrefix(err.Error(), "nats: ")
}
