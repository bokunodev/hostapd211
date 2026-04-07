// Package hostapd211 provides a Go client for the hostapd control interface.
//
// It allows for sending commands and monitoring events from a hostapd instance.
//
// Example usage:
//
//	client := hostapd211.NewClient("/var/run/hostapd/wlan0")
//	reply, err := client.Ping(ctx)
package hostapd211

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const HOSTAPD_VERSION = "2.11"

type Command string

const (
	PING                        Command = "PING"
	RELOG                       Command = "RELOG"
	CLOSE_LOG                   Command = "CLOSE_LOG"
	NOTE                        Command = "NOTE"
	STATUS                      Command = "STATUS"
	STATUS_DRIVER               Command = "STATUS-DRIVER"
	MIB                         Command = "MIB"
	STA_FIRST                   Command = "STA-FIRST"
	STA                         Command = "STA"
	STA_NEXT                    Command = "STA-NEXT"
	ATTACH                      Command = "ATTACH"
	DETACH                      Command = "DETACH"
	LEVEL                       Command = "LEVEL"
	NEW_STA                     Command = "NEW_STA"
	DEAUTHENTICATE              Command = "DEAUTHENTICATE"
	DISASSOCIATE                Command = "DISASSOCIATE"
	SIGNATURE                   Command = "SIGNATURE"
	POLL_STA                    Command = "POLL_STA"
	STOP_AP                     Command = "STOP_AP"
	SA_QUERY                    Command = "SA_QUERY"
	WPS_PIN                     Command = "WPS_PIN"
	WPS_CHECK_PIN               Command = "WPS_CHECK_PIN"
	WPS_PBC                     Command = "WPS_PBC"
	WPS_CANCEL                  Command = "WPS_CANCEL"
	WPS_AP_PIN                  Command = "WPS_AP_PIN"
	WPS_CONFIG                  Command = "WPS_CONFIG"
	WPS_GET_STATUS              Command = "WPS_GET_STATUS"
	WPS_NFC_TAG_READ            Command = "WPS_NFC_TAG_READ"
	WPS_NFC_CONFIG_TOKEN        Command = "WPS_NFC_CONFIG_TOKEN"
	WPS_NFC_TOKEN               Command = "WPS_NFC_TOKEN"
	NFC_GET_HANDOVER_SEL        Command = "NFC_GET_HANDOVER_SEL"
	NFC_REPORT_HANDOVER         Command = "NFC_REPORT_HANDOVER"
	SET_QOS_MAP_SET             Command = "SET_QOS_MAP_SET"
	SEND_QOS_MAP_CONF           Command = "SEND_QOS_MAP_CONF"
	HS20_WNM_NOTIF              Command = "HS20_WNM_NOTIF"
	HS20_DEAUTH_REQ             Command = "HS20_DEAUTH_REQ"
	DISASSOC_IMMINENT           Command = "DISASSOC_IMMINENT"
	ESS_DISASSOC                Command = "ESS_DISASSOC"
	BSS_TM_REQ                  Command = "BSS_TM_REQ"
	COLOC_INTF_REQ              Command = "COLOC_INTF_REQ"
	GET_CONFIG                  Command = "GET_CONFIG"
	SET                         Command = "SET"
	GET                         Command = "GET"
	ENABLE                      Command = "ENABLE"
	RELOAD_WPA_PSK              Command = "RELOAD_WPA_PSK"
	GET_RXKHS                   Command = "GET_RXKHS"
	RELOAD_RXKHS                Command = "RELOAD_RXKHS"
	RELOAD_BSS                  Command = "RELOAD_BSS"
	RELOAD_CONFIG               Command = "RELOAD_CONFIG"
	RELOAD                      Command = "RELOAD"
	DISABLE                     Command = "DISABLE"
	UPDATE_BEACON               Command = "UPDATE_BEACON"
	CHAN_SWITCH                 Command = "CHAN_SWITCH"
	COLOR_CHANGE                Command = "COLOR_CHANGE"
	NOTIFY_CW_CHANGE            Command = "NOTIFY_CW_CHANGE"
	VENDOR                      Command = "VENDOR"
	ERP_FLUSH                   Command = "ERP_FLUSH"
	EAPOL_REAUTH                Command = "EAPOL_REAUTH"
	EAPOL_SET                   Command = "EAPOL_SET"
	LOG_LEVEL                   Command = "LOG_LEVEL"
	TRACK_STA_LIST              Command = "TRACK_STA_LIST"
	PMKSA                       Command = "PMKSA"
	PMKSA_FLUSH                 Command = "PMKSA_FLUSH"
	PMKSA_ADD                   Command = "PMKSA_ADD"
	SET_NEIGHBOR                Command = "SET_NEIGHBOR"
	SHOW_NEIGHBOR               Command = "SHOW_NEIGHBOR"
	REMOVE_NEIGHBOR             Command = "REMOVE_NEIGHBOR"
	REQ_LCI                     Command = "REQ_LCI"
	REQ_RANGE                   Command = "REQ_RANGE"
	REQ_BEACON                  Command = "REQ_BEACON"
	REQ_LINK_MEASUREMENT        Command = "REQ_LINK_MEASUREMENT"
	TERMINATE                   Command = "TERMINATE"
	ACCEPT_ACL                  Command = "ACCEPT_ACL"
	DENY_ACL                    Command = "DENY_ACL"
	DPP_QR_CODE                 Command = "DPP_QR_CODE"
	DPP_NFC_URI                 Command = "DPP_NFC_URI"
	DPP_NFC_HANDOVER_REQ        Command = "DPP_NFC_HANDOVER_REQ"
	DPP_NFC_HANDOVER_SEL        Command = "DPP_NFC_HANDOVER_SEL"
	DPP_BOOTSTRAP_GEN           Command = "DPP_BOOTSTRAP_GEN"
	DPP_BOOTSTRAP_REMOVE        Command = "DPP_BOOTSTRAP_REMOVE"
	DPP_BOOTSTRAP_GET_URI       Command = "DPP_BOOTSTRAP_GET_URI"
	DPP_BOOTSTRAP_INFO          Command = "DPP_BOOTSTRAP_INFO"
	DPP_BOOTSTRAP_SET           Command = "DPP_BOOTSTRAP_SET"
	DPP_AUTH_INIT               Command = "DPP_AUTH_INIT"
	DPP_LISTEN                  Command = "DPP_LISTEN"
	DPP_STOP_LISTEN             Command = "DPP_STOP_LISTEN"
	DPP_CONFIGURATOR_ADD        Command = "DPP_CONFIGURATOR_ADD"
	DPP_CONFIGURATOR_SET        Command = "DPP_CONFIGURATOR_SET"
	DPP_CONFIGURATOR_REMOVE     Command = "DPP_CONFIGURATOR_REMOVE"
	DPP_CONFIGURATOR_SIGN       Command = "DPP_CONFIGURATOR_SIGN"
	DPP_CONFIGURATOR_GET_KEY    Command = "DPP_CONFIGURATOR_GET_KEY"
	DPP_PKEX_ADD                Command = "DPP_PKEX_ADD"
	DPP_PKEX_REMOVE             Command = "DPP_PKEX_REMOVE"
	DPP_CONTROLLER_START        Command = "DPP_CONTROLLER_START"
	DPP_CONTROLLER_STOP         Command = "DPP_CONTROLLER_STOP"
	DPP_CHIRP                   Command = "DPP_CHIRP"
	DPP_STOP_CHIRP              Command = "DPP_STOP_CHIRP"
	DPP_RELAY_ADD_CONTROLLER    Command = "DPP_RELAY_ADD_CONTROLLER"
	DPP_RELAY_REMOVE_CONTROLLER Command = "DPP_RELAY_REMOVE_CONTROLLER"
	NAN_PUBLISH                 Command = "NAN_PUBLISH"
	NAN_CANCEL_PUBLISH          Command = "NAN_CANCEL_PUBLISH"
	NAN_UPDATE_PUBLISH          Command = "NAN_UPDATE_PUBLISH"
	NAN_SUBSCRIBE               Command = "NAN_SUBSCRIBE"
	NAN_CANCEL_SUBSCRIBE        Command = "NAN_CANCEL_SUBSCRIBE"
	NAN_TRANSMIT                Command = "NAN_TRANSMIT"
	GET_CAPABILITY              Command = "GET_CAPABILITY"
	PTKSA_CACHE_LIST            Command = "PTKSA_CACHE_LIST"
	DRIVER                      Command = "DRIVER"
	ENABLE_MLD                  Command = "ENABLE_MLD"
	DISABLE_MLD                 Command = "DISABLE_MLD"
)

type Client struct {
	remote_addr *net.UnixAddr
	local_addr  *net.UnixAddr
	conn        *net.UnixConn
}

// [FindRemoteSocket] scans the socket directory and returns the first socket file if socket_name is empty
// and returns an error if no socket file is found
func FindRemoteSocket(socket_dir, socket_name string) (string, error) {
	if socket_name == "" {
		entries, err := os.ReadDir(socket_dir)
		if err != nil {
			return "", err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if entry.Type() == os.ModeSocket {
				socket_name = entry.Name()
			}
		}

		if socket_name == "" {
			return "", fmt.Errorf("no socket found")
		}
	}

	return filepath.Join(socket_dir, socket_name), nil
}

func NewClient(remote_socket string) (*Client, error) {
	remote_addr, err := net.ResolveUnixAddr("unixgram", remote_socket)
	if err != nil {
		return nil, err
	}

	local_socket := fmt.Sprintf("/tmp/hostapd-client-%d-%d.sock", os.Getpid(), time.Now().UnixNano())
	os.Remove(local_socket)

	local_addr, err := net.ResolveUnixAddr("unixgram", local_socket)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUnix("unixgram", local_addr, remote_addr)
	if err != nil {
		return nil, err
	}

	return &Client{remote_addr: remote_addr, local_addr: local_addr, conn: conn}, nil
}

func (c *Client) Close() error {
	defer os.Remove(c.local_addr.Name)

	return c.conn.Close()
}

// Do sends command to hostapd
//
// read and write deadlines are set to ctx.Deadline()
func (c *Client) Do(ctx context.Context, cmd Command, args ...string) (string, error) {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(string(cmd))
	for _, arg := range args {
		buffer.WriteByte(' ')
		buffer.WriteString(arg)
	}

	if deadline, ok := ctx.Deadline(); ok {
		err := c.conn.SetDeadline(deadline)
		if err != nil {
			return "", err
		}
	}

	_, err := c.conn.Write(buffer.Bytes())
	if err != nil {
		return "", err
	}

	reply := [4 << 10]byte{}
	n, err := c.conn.Read(reply[:])
	if err != nil {
		return "", err
	}

	return string(reply[:n]), nil
}

// Ping sends [PING] command to hostapd
func (c *Client) Ping(ctx context.Context) (string, error) {
	return c.Do(ctx, PING)
}

// Relog sends [RELOG] command to hostapd
func (c *Client) Relog(ctx context.Context) (string, error) {
	return c.Do(ctx, RELOG)
}

// CloseLog sends [CLOSE_LOG] command to hostapd
func (c *Client) CloseLog(ctx context.Context) (string, error) {
	return c.Do(ctx, CLOSE_LOG)
}

// Note sends [NOTE] command to hostapd
func (c *Client) Note(ctx context.Context, msg string) (string, error) {
	return c.Do(ctx, NOTE, msg)
}

// Status sends [STATUS] command to hostapd
func (c *Client) Status(ctx context.Context) (string, error) {
	return c.Do(ctx, STATUS)
}

// StatusDriver sends [STATUS_DRIVER] to hostapd
func (c *Client) StatusDriver(ctx context.Context) (string, error) {
	return c.Do(ctx, STATUS_DRIVER)
}

// MIB sends [MIB] command to hostapd
func (c *Client) MIB(ctx context.Context, radius_server bool) (string, error) {
	args := []string{}
	if radius_server {
		args = append(args, "radius_server")
	}

	return c.Do(ctx, MIB, args...)
}

// STAFirst sends [STA_FIRST] command to hostapd
func (c *Client) STAFirst(ctx context.Context) (string, error) {
	return c.Do(ctx, STA_FIRST)
}

// STA sends [STA] command to hostapd
func (c *Client) STA(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, STA, addr.String())
}

// STANext sends [STA_NEXT] command to hostapd
func (c *Client) STANext(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, STA_NEXT, addr.String())
}

var ErrStop = errors.New("stop")

// Attach sends [ATTACH] command to hostapd
//
// Attach create a separate UNIX socket for listening on events
// and wait for events until callback returns an error or [ErrStop].
// if callback returns [ErrStop], [Client.Attach] returns nil.
func (c *Client) Attach(
	ctx context.Context,
	read_timeout time.Duration,
	callback func(ctx context.Context, msg string) error,
) error {
	local_socket := fmt.Sprintf("/tmp/hostapd-client-%d-%d.sock", os.Getpid(), time.Now().UnixNano())
	// remove existing socket file if exists
	// to prevent stale socket file
	os.Remove(local_socket)

	// create a new local unix socket exclusively for listening on events
	local_addr, err := net.ResolveUnixAddr("unixgram", local_socket)
	if err != nil {
		return err
	}

	conn, err := net.DialUnix("unixgram", local_addr, c.remote_addr)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(ATTACH))
	if err != nil {
		return err
	}

	msgbuf := [4 << 10]byte{}
	n, err := conn.Read(msgbuf[:])
	if err != nil {
		return err
	}

	if string(msgbuf[:n]) != "OK\n" {
		return errors.New("FAIL")
	}

	defer func() {
		// detach on exit, no need to check the reply
		// because the connection will be abandoned anyway
		conn.Write([]byte(DETACH))
		conn.Close()
		os.Remove(local_addr.Name)
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if read_timeout > 0 {
			// set read deadline to prevent blocking forever
			err = conn.SetReadDeadline(time.Now().Add(read_timeout))
			if err != nil {
				return err
			}
		}

		n, err := conn.Read(msgbuf[:])
		if err != nil {
			// if read timed out, try again (allows the ctx.Done() to be checked)
			if nerr, ok := errors.AsType[net.Error](err); ok && nerr.Timeout() {
				continue loop
			}

			return err
		}

		err = callback(ctx, string(bytes.TrimSpace(msgbuf[:n])))
		if err != nil {
			if errors.Is(err, ErrStop) {
				return nil
			}

			return err
		}
	}
}

// Detach sends [DETACH] command to hostapd
func (c *Client) Detach(ctx context.Context) (string, error) {
	return c.Do(ctx, DETACH)
}

// Level sends [LEVEL] command to hostapd
func (c *Client) Level(ctx context.Context, level int) (string, error) {
	return c.Do(ctx, LEVEL, strconv.Itoa(level))
}

// NewSTA sends [NEW_STA] command to hostapd
func (c *Client) NewSTA(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, NEW_STA, addr.String())
}

// Deauthenticate sends [DEAUTHENTICATE] command to hostapd,
//
// reason and silent are optional
func (c *Client) Deauthenticate(
	ctx context.Context,
	addr net.HardwareAddr,
	reason string,
	silent bool,
) (string, error) {
	args := []string{}

	if reason != "" {
		args = append(args, fmt.Sprintf("reason=%s", reason))
	}

	if silent {
		args = append(args, "tx=0")
	}

	return c.Do(ctx, DEAUTHENTICATE, args...)
}

// Disassociate sends [DISASSOCIATE] command to hostapd,
//
// reason and silent are optional
func (c *Client) Disassociate(
	ctx context.Context,
	addr net.HardwareAddr,
	reason string,
	silent bool,
) (string, error) {
	args := []string{}

	if reason != "" {
		args = append(args, fmt.Sprintf("reason=%s", reason))
	}

	if silent {
		args = append(args, "tx=0")
	}

	return c.Do(ctx, DISASSOCIATE, args...)
}

// Signature sends [SIGNATURE] command to hostapd
func (c *Client) Signature(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, SIGNATURE, addr.String())
}

// PollSTA sends [POLL_STA] command to hostapd
func (c *Client) PollSTA(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, POLL_STA, addr.String())
}

// StopAP sends [STOP_AP] command to hostapd
func (c *Client) StopAP(ctx context.Context) (string, error) {
	return c.Do(ctx, STOP_AP)
}

// SAQuery sends [SA_QUERY] command to hostapd
func (c *Client) SAQuery(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, SA_QUERY, addr.String())
}

// WPSPin sends [WPS_PIN] command to hostapd
//
// uuid must be either `any` or a valid UUID
func (c *Client) WPSPin(
	ctx context.Context,
	pin string,
	timeout time.Duration,
	addr net.HardwareAddr,
	uuid string,
) (string, error) {
	return c.Do(ctx, WPS_PIN, pin, strconv.Itoa(int(timeout.Seconds())), addr.String(), uuid)
}

// WPSCheckPin sends [WPS_CHECK_PIN] command to hostapd
func (c *Client) WPSCheckPin(ctx context.Context, pin string) (string, error) {
	return c.Do(ctx, WPS_CHECK_PIN, pin)
}

// WPSButton sends [WPS_PBC] command to hostapd
func (c *Client) WPSButton(ctx context.Context) (string, error) {
	return c.Do(ctx, WPS_PBC)
}

// WPSCancel sends [WPS_CANCEL] command to hostapd
func (c *Client) WPSCancel(ctx context.Context) (string, error) {
	return c.Do(ctx, WPS_CANCEL)
}

type WPSAPPinCmd string

const (
	WPSAPPinCmdDisable WPSAPPinCmd = "disable"
	WPSAPPinCmdRandom  WPSAPPinCmd = "random"
	WPSAPPinCmdGet     WPSAPPinCmd = "get"
	WPSAPPinCmdSet     WPSAPPinCmd = "set"
)

// WPSAPPin sends [WPS_AP_PIN] command to hostapd
//
// pin and timeout are ignored if cmd is not [WPSAPPinCmdSet]
func (c *Client) WPSAPPin(
	ctx context.Context,
	cmd WPSAPPinCmd,
	pin string,
	timeout time.Duration,
) (string, error) {
	args := []string{}
	args = append(args, string(cmd))

	if cmd == WPSAPPinCmdGet {
		if pin == "" || timeout <= 0 {
			return "", fmt.Errorf("`pin` and `timeout` are required for `WPSAPPinSet`")
		}

		args = append(args, pin, strconv.Itoa(int(timeout.Seconds())))
	}

	return c.Do(ctx, WPS_AP_PIN, args...)
}

type WPSConfigAuth string

const (
	WPSConfigAuthOpen    WPSConfigAuth = "OPEN"
	WPSConfigAuthWPAPSK  WPSConfigAuth = "WPAPSK"
	WPSConfigAuthWPA2PSK WPSConfigAuth = "WPA2PSK"
)

type WPSConfigEncr string

const (
	WPSConfigEncrNone WPSConfigEncr = "NONE"
	WPSConfigEncrTKIP WPSConfigEncr = "TKIP"
	WPSConfigEncrCCMP WPSConfigEncr = "CCMP"
)

// WPSConfig sends [WPS_CONFIG] command to hostapd
//
// key is require if encr is not [WPSConfigEncrNone]
func (c *Client) WPSConfig(
	ctx context.Context,
	ssid string,
	auth WPSConfigAuth,
	encr WPSConfigEncr,
	key string,
) (string, error) {
	args := []string{}

	args = append(args, ssid)
	args = append(args, string(auth))
	args = append(args, string(encr))

	if encr != WPSConfigEncrNone {
		args = append(args, key)
	}

	return c.Do(ctx, WPS_CONFIG, args...)
}

// WPSGetStatus sends [WPS_GET_STATUS] command to hostapd
func (c *Client) WPSGetStatus(ctx context.Context) (string, error) {
	return c.Do(ctx, WPS_GET_STATUS)
}

// WPSNFCTagRead sends [WPS_NFC_TAG_READ] command to hostapd
//
// data is NFC tag data in hexadecimal format
func (c *Client) WPSNFCTagRead(ctx context.Context, data string) (string, error) {
	return c.Do(ctx, WPS_NFC_TAG_READ, data)
}

type WPSNFCConfigTokenCmd string

const (
	WPSNFCConfigTokenCmdWPS  WPSNFCConfigTokenCmd = "WPS"
	WPSNFCConfigTokenCmdNDEF WPSNFCConfigTokenCmd = "NDEF"
)

// WPSNFCConfigToken sends [WPS_NFC_CONFIG_TOKEN] command to hostapd
func (c *Client) WPSNFCConfigToken(ctx context.Context, ndef WPSNFCConfigTokenCmd) (string, error) {
	return c.Do(ctx, WPS_NFC_CONFIG_TOKEN, string(ndef))
}

type WPSNFCTokenCmd string

const (
	WPSNFCTokenCmdWPS     WPSNFCTokenCmd = "WPS"
	WPSNFCTokenCmdNDEF    WPSNFCTokenCmd = "NDEF"
	WPSNFCTokenCmdEnable  WPSNFCTokenCmd = "enable"
	WPSNFCTokenCmdDisable WPSNFCTokenCmd = "disable"
)

// WPSNFCToken sends [WPS_NFC_TOKEN] command to hostapd
func (c *Client) WPSNFCToken(ctx context.Context, cmd WPSNFCTokenCmd) (string, error) {
	return c.Do(ctx, WPS_NFC_TOKEN, string(cmd))
}

type NFCGetHandoverSelCmd string

const (
	NFCGetHandoverSelCmdWPS  NFCGetHandoverSelCmd = "WPS"
	NFCGetHandoverSelCmdNDEF NFCGetHandoverSelCmd = "NDEF"
)

// NFC_GET_HANDOVER_SEL
func (c *Client) NFCGetHandoverSel(
	ctx context.Context,
	cmd NFCGetHandoverSelCmd,
	wps_cr bool,
) (string, error) {
	args := []string{}
	args = append(args, string(cmd))

	if wps_cr {
		args = append(args, "WPS-CR")
	}

	return c.Do(ctx, NFC_GET_HANDOVER_SEL, args...)
}

// NFCReportHandover sends [NFC_REPORT_HANDOVER] command to hostapd
//
// role and typ must be either `RESP` and `WPS` respectively
//
// req_hex and sel_hex are hexadecimal strings of NFC handover request and
// selection data respectively
func (c *Client) NFCReportHandover(ctx context.Context, role, typ, req_hex, sel_hex string) (string, error) {
	if role != "RESP" {
		return "", fmt.Errorf("`role` must be `RESP`")
	}

	if typ != "WPS" {
		return "", fmt.Errorf("`typ` must be `WPS`")
	}

	return c.Do(ctx, NFC_REPORT_HANDOVER, role, typ, req_hex, sel_hex)
}

// SetQOSMapSet sends [SET_QOS_MAP_SET] command to hostapd
//
// sets must be comma separated numbers between 0-255 inclusive
func (c *Client) SetQOSMapSet(ctx context.Context, sets string) (string, error) {
	return c.Do(ctx, SET_QOS_MAP_SET, sets)
}

// SendQOSMapConf sends [SEND_QOS_MAP_CONF] command to hostapd
func (c *Client) SendQOSMapConf(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, SEND_QOS_MAP_CONF, addr.String())
}

// HS20WNMNotif sends [HS20_WNM_NOTIF] command to hostapd
func (c *Client) HS20WNMNotif(ctx context.Context, addr net.HardwareAddr, url string) (string, error) {
	return c.Do(ctx, HS20_WNM_NOTIF, addr.String(), url)
}

// HS20DeauthReq sends [HS20_DEAUTH_REQ] command to hostapd
func (c *Client) HS20DeauthReq(
	ctx context.Context,
	addr net.HardwareAddr,
	code bool,
	reauth_delay time.Duration,
	url string,
) (string, error) {
	args := []string{}
	args = append(args, addr.String())

	if code {
		args = append(args, "1")
	} else {
		args = append(args, "0")
	}

	args = append(args, strconv.Itoa(int(reauth_delay.Seconds())))
	args = append(args, url)

	return c.Do(ctx, HS20_DEAUTH_REQ, args...)
}

// DisassocImminent sends [DISASSOC_IMMINENT] command to hostapd
func (c *Client) DisassocImminent(ctx context.Context, addr net.HardwareAddr, timer time.Duration) (string, error) {
	return c.Do(ctx, DISASSOC_IMMINENT, addr.String(), strconv.Itoa(int(timer.Seconds())))
}

// ESSDisassoc sends [ESS_DISASSOC] command to hostapd
func (c *Client) ESSDisassoc(ctx context.Context, addr net.HardwareAddr, timer time.Duration) (string, error) {
	return c.Do(ctx, ESS_DISASSOC, addr.String(), strconv.Itoa(int(timer.Seconds())))
}

//	BSSTMReq sends [BSS_TM_REQ] command to hostapd
//
// neighbor format is `<bssid:mac address>,<bssid information:32-bit int>,<operating class 8-bit int>,<channel number 8-bit int>,<phy type 8-bit int>,<optional subelements hex string>`
// mbo format is `<reason:int>:<reassoc_delay:int>:<cell_pref:int>`
func (c *Client) BSSTMReq(
	ctx context.Context,
	addr net.HardwareAddr,
	disassoc_timer, valid_int, dialog_token int,
	bss_term, url, neighbor, mbo string,
) (string, error) {
	args := []string{}
	args = append(args, addr.String())
	args = append(args, fmt.Sprintf("disassoc_timer=%d", disassoc_timer))
	args = append(args, fmt.Sprintf("valid_int=%d", valid_int))
	args = append(args, fmt.Sprintf("dialog_token=%d", dialog_token))
	args = append(args, fmt.Sprintf("bss_term=%s", bss_term))
	args = append(args, fmt.Sprintf("url=%s", url))
	args = append(args, fmt.Sprintf("neighbor=%s", neighbor))
	args = append(args, fmt.Sprintf("mbo=%s", mbo))

	return c.Do(ctx, BSS_TM_REQ, args...)
}

// ColocIntfReq sends [COLOC_INTF_REQ] command to hostapd
func (c *Client) ColocIntfReq(ctx context.Context, addr net.HardwareAddr, auto_report, timeout int) (string, error) {
	return c.Do(ctx, COLOC_INTF_REQ, addr.String(), strconv.Itoa(auto_report), strconv.Itoa(timeout))
}

// GetConfig sends [GET_CONFIG] command to hostapd
func (c *Client) GetConfig(ctx context.Context) (string, error) {
	return c.Do(ctx, GET_CONFIG)
}

// Set sends [SET] command to hostapd
func (c *Client) Set(ctx context.Context, key, value string) (string, error) {
	return c.Do(ctx, SET, key, value)
}

// Get sends [GET] command to hostapd
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Do(ctx, GET, key)
}

// Enable sends [ENABLE] command to hostapd
func (c *Client) Enable(ctx context.Context) (string, error) {
	return c.Do(ctx, ENABLE)
}

// ReloadWPAPSK sends [RELOAD_WPA_PSK] command to hostapd
func (c *Client) ReloadWPAPSK(ctx context.Context) (string, error) {
	return c.Do(ctx, RELOAD_WPA_PSK)
}

// GetRXKHS sends [GET_RXKHS] command to hostapd
func (c *Client) GetRXKHS(ctx context.Context) (string, error) {
	return c.Do(ctx, GET_RXKHS)
}

// ReloadRXKHS sends [RELOAD_RXKHS] command to hostapd
func (c *Client) ReloadRXKHS(ctx context.Context) (string, error) {
	return c.Do(ctx, RELOAD_RXKHS)
}

// ReloadBSS sends [RELOAD_BSS] command to hostapd
func (c *Client) ReloadBSS(ctx context.Context) (string, error) {
	return c.Do(ctx, RELOAD_BSS)
}

// ReloadConfig sends [RELOAD_CONFIG] command to hostapd
func (c *Client) ReloadConfig(ctx context.Context) (string, error) {
	return c.Do(ctx, RELOAD_CONFIG)
}

// Reload sends [RELOAD] command to hostapd
func (c *Client) Reload(ctx context.Context) (string, error) {
	return c.Do(ctx, RELOAD)
}

// Disable sends [DISABLE] command to hostapd
func (c *Client) Disable(ctx context.Context) (string, error) {
	return c.Do(ctx, DISABLE)
}

// UpdateBeacon sends [UPDATE_BEACON] command to hostapd
func (c *Client) UpdateBeacon(ctx context.Context) (string, error) {
	return c.Do(ctx, UPDATE_BEACON)
}

// ChanSwitch sends [CHAN_SWITCH] command to hostapd
//
//	center_freq1       => The center frequency of the first segment (for 80/160 MHz).
//	center_freq2       => The center frequency of the second segment (for 80+80 MHz).
//	bandwidth          => The channel width (e.g., 20, 40, 80, 160).
//	sec_channel_offset => Secondary channel offset for HT40 (+1 or -1).
//	punct_bitmap       => Preamble puncturing bitmap for wider channels.
//
//	ht      => Enables High Throughput (802.11n).
//	vht     => Enables Very High Throughput (802.11ac).
//	he      => Enables High Efficiency (802.11ax/Wi-Fi 6).
//	eht     => Enables Extremely High Throughput (802.11be/Wi-Fi 7).
//	blocktx => Instructs the AP to block transmissions on the current channel until the
func (c *Client) ChanSwitch(
	ctx context.Context,
	cs_count, freq int,
	center_freq1, center_freq2, bandwidth, sec_channel_offset, punct_bitmap string,
	ht, vht, he, eht, blocktx bool,
) (string, error) {
	args := []string{}
	args = append(args, fmt.Sprintf("cs_count=%d", cs_count))
	args = append(args, fmt.Sprintf("freq=%d", freq))
	args = append(args, fmt.Sprintf("center_freq1=%s", center_freq1))
	args = append(args, fmt.Sprintf("center_freq2=%s", center_freq2))
	args = append(args, fmt.Sprintf("bandwidth=%s", bandwidth))
	args = append(args, fmt.Sprintf("sec_channel_offset=%s", sec_channel_offset))
	args = append(args, fmt.Sprintf("punct_bitmap=%s", punct_bitmap))

	if ht {
		args = append(args, "ht")
	}

	if vht {
		args = append(args, "vht")
	}

	if he {
		args = append(args, "he")
	}

	if eht {
		args = append(args, "eht")
	}

	if blocktx {
		args = append(args, "blocktx")
	}

	return c.Do(ctx, CHAN_SWITCH, args...)
}

//	ColorChange sends [COLOR_CHANGE] command to hostapd
//
// color must be between 0-63 inclusive
func (c *Client) ColorChange(ctx context.Context, color int) (string, error) {
	return c.Do(ctx, COLOR_CHANGE, strconv.Itoa(color))
}

// NotifyCWChange sends [NOTIFY_CW_CHANGE] command to hostapd
//
// cw must be either 0, 1, 2, or 3
func (c *Client) NotifyCWChange(ctx context.Context, cw int) (string, error) {
	return c.Do(ctx, NOTIFY_CW_CHANGE, strconv.Itoa(cw))
}

// Vendor sends [VENDOR] command to hostapd
//
// data is optional hexadecimal string of vendor specific data
// nested is ignored if data is not provided
func (c *Client) Vendor(
	ctx context.Context,
	vendor_id string,
	subcommand string,
	data string,
	nested bool,
) (string, error) {
	args := []string{}
	args = append(args, vendor_id)
	args = append(args, subcommand)

	if data != "" {
		args = append(args, data)

		if nested {
			args = append(args, "nested=1")
		} else {
			args = append(args, "nested=0")
		}
	}

	return c.Do(ctx, VENDOR, args...)
}

// ERPFlush sends [ERP_FLUSH] command to hostapd
func (c *Client) ERPFlush(ctx context.Context) (string, error) {
	return c.Do(ctx, ERP_FLUSH)
}

// EAPOLReauth sends [EAPOL_REAUTH] command to hostapd
func (c *Client) EAPOLReauth(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, EAPOL_REAUTH, addr.String())
}

// EAPOLSet sends [EAPOL_SET] command to hostapd
func (c *Client) EAPOLSet(ctx context.Context, addr net.HardwareAddr, param, value string) (string, error) {
	return c.Do(ctx, EAPOL_SET, addr.String(), param, value)
}

type LogLevelCmd string

const (
	LogLevelExcessive = "EXCESSIVE"
	LogLevelMsgDump   = "MSGDUMP"
	LogLevelDebug     = "DEBUG"
	LogLevelInfo      = "INFO"
	LogLevelWarning   = "WARNING"
	LogLevelError     = "ERROR"
)

// LogLevel sends [LOG_LEVEL] command to hostapd
func (c *Client) LogLevel(ctx context.Context, level LogLevelCmd) (string, error) {
	return c.Do(ctx, LOG_LEVEL, string(level))
}

// TrackSTAList sends [TRACK_STA_LIST] command to hostapd
func (c *Client) TrackSTAList(ctx context.Context) (string, error) {
	return c.Do(ctx, TRACK_STA_LIST)
}

// PMKSA sends [PMKSA] command to hostapd
func (c *Client) PMKSA(ctx context.Context) (string, error) {
	return c.Do(ctx, PMKSA)
}

// PMKSAFlush sends [PMKSA_FLUSH] command to hostapd
func (c *Client) PMKSAFlush(ctx context.Context) (string, error) {
	return c.Do(ctx, PMKSA_FLUSH)
}

// PMKSAAdd sends [PMKSA_ADD] command to hostapd
func (c *Client) PMKSAAdd(ctx context.Context,
	addr net.HardwareAddr,
	pmkid string,
	pmk string,
	expiration time.Duration,
	akmp int,
) (string, error) {
	return c.Do(ctx, PMKSA_ADD, addr.String(), pmkid, pmk, strconv.Itoa(int(expiration.Seconds())), strconv.Itoa(akmp))
}

// SetNeighbor sends [SET_NEIGHBOR] command to hostapd
//
// lci, civic, bss_parameter and stat are optional
func (c *Client) SetNeighbor(
	ctx context.Context,
	bssid net.HardwareAddr,
	ssid string,
	nr string,
	lci string,
	civic string,
	bss_parameter int,
	stat bool,
) (string, error) {
	args := []string{}
	args = append(args, bssid.String())
	args = append(args, fmt.Sprintf("ssid=%s", ssid))
	args = append(args, fmt.Sprintf("nr=%s", nr))

	if lci != "" {
		args = append(args, fmt.Sprintf("lci=%s", lci))
	}

	if civic != "" {
		args = append(args, fmt.Sprintf("civic=%s", civic))
	}

	if bss_parameter > 0 {
		args = append(args, fmt.Sprintf("bss_parameter=%d", bss_parameter))
	}

	if stat {
		args = append(args, "stat")
	}

	return c.Do(ctx, SET_NEIGHBOR, args...)
}

// ShowNeighbor sends [SHOW_NEIGHBOR] command to hostapd
func (c *Client) ShowNeighbor(ctx context.Context) (string, error) {
	return c.Do(ctx, SHOW_NEIGHBOR)
}

// RemoveNeighbor sends [REMOVE_NEIGHBOR] command to hostapd
func (c *Client) RemoveNeighbor(ctx context.Context, bssid net.HardwareAddr, ssid string) (string, error) {
	return c.Do(ctx, REMOVE_NEIGHBOR, bssid.String(), fmt.Sprintf("ssid=%s", ssid))
}

// ReqLCI sends [REQ_LCI] command to hostapd
func (c *Client) ReqLCI(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, REQ_LCI, addr.String())
}

// ReqRange sends [REQ_RANGE] command to hostapd
//
// responders is a list of one or more MAC addresses representing the FTM responders
func (c *Client) ReqRange(ctx context.Context,
	addr net.HardwareAddr,
	random_interval int,
	min_ap int,
	responders string,
) (string, error) {
	args := []string{}
	args = append(args, addr.String())
	args = append(args, strconv.Itoa(random_interval))
	args = append(args, strconv.Itoa(min_ap))
	args = append(args, responders)

	return c.Do(ctx, REQ_RANGE, args...)
}

// ReqBeacon sends [REQ_BEACON] command to hostapd
//
// req_mode and req_data are optional hex-encoded strings
func (c *Client) ReqBeacon(ctx context.Context, addr net.HardwareAddr, req_mode string, req_data string) (string, error) {
	args := []string{}
	args = append(args, addr.String())
	if req_mode != "" {
		args = append(args, fmt.Sprintf("req_mode=%s", req_mode))
	}

	args = append(args, req_data)

	return c.Do(ctx, REQ_BEACON, args...)
}

// ReqLinkMeasurement sends [REQ_LINK_MEASUREMENT] command to hostapd
func (c *Client) ReqLinkMeasurement(ctx context.Context, addr net.HardwareAddr) (string, error) {
	return c.Do(ctx, REQ_LINK_MEASUREMENT, addr.String())
}

// Terminate sends [TERMINATE] command to hostapd
func (c *Client) Terminate(ctx context.Context) (string, error) {
	return c.Do(ctx, TERMINATE)
}

type AcceptACLCmd string

const (
	AcceptAclAddMac AcceptACLCmd = "ADD_MAC"
	AcceptAclDelMac AcceptACLCmd = "DEL_MAC"
	AcceptAclShow   AcceptACLCmd = "SHOW"
	AcceptAclClear  AcceptACLCmd = "CLEAR"
)

// AcceptACL sends [ACCEPT_ACL] command to hostapd
//
// addr is only required for ADD_MAC and DEL_MAC commands and vlan is optional for ADD_MAC command
func (c *Client) AcceptACL(ctx context.Context, cmd AcceptACLCmd, addr net.HardwareAddr, vlan int) (string, error) {
	args := []string{}
	args = append(args, string(cmd))

	switch cmd {
	case AcceptAclAddMac:
		args = append(args, addr.String())
		if vlan > 0 {
			args = append(args, fmt.Sprintf("VLAND_ID=%d", vlan))
		}

	case AcceptAclDelMac:
		args = append(args, addr.String())
	default:
	}

	return c.Do(ctx, ACCEPT_ACL, args...)
}

type DenyACLCmd string

const (
	DenyAclAddMac DenyACLCmd = "ADD_MAC"
	DenyAclDelMac DenyACLCmd = "DEL_MAC"
	DenyAclShow   DenyACLCmd = "SHOW"
	DenyAclClear  DenyACLCmd = "CLEAR"
)

// DenyACL sends [DENY_ACL] command to hostapd
//
// addr is only required for ADD_MAC and DEL_MAC commands and vlan is optional for ADD_MAC command
func (c *Client) DenyACL(ctx context.Context, cmd DenyACLCmd, addr net.HardwareAddr, vlan int) (string, error) {
	args := []string{}
	args = append(args, string(cmd))

	switch cmd {
	case DenyAclAddMac:
		args = append(args, addr.String())
		if vlan > 0 {
			args = append(args, fmt.Sprintf("VLAND_ID=%d", vlan))
		}

	case DenyAclDelMac:
		args = append(args, addr.String())
	default:
	}

	return c.Do(ctx, DENY_ACL, args...)
}

// DPPQRCode sends [DPP_QR_CODE] command to hostapd
//
// uri is DPP data uri from the QR code
func (c *Client) DPPQRCode(ctx context.Context, uri string) (string, error) {
	return c.Do(ctx, DPP_QR_CODE, uri)
}

// DPPNFCURI sends [DPP_NFC_URI] command to hostapd
//
// uri is NFC tag data uri
func (c *Client) DPPNFCURI(ctx context.Context, uri string) (string, error) {
	return c.Do(ctx, DPP_NFC_URI, uri)
}

// DPPNFCHandoverReq sends [DPP_NFC_HANDOVER_REQ] command to hostapd
//
// own is DPP record id and uri is NFC tag data uri
func (c *Client) DPPNFCHandoverReq(ctx context.Context, own int, uri string) (string, error) {
	return c.Do(ctx, DPP_NFC_HANDOVER_REQ, fmt.Sprintf("own=%d", own), fmt.Sprintf("uri=%s", uri))
}

// DPPNFCHandoverSel sends [DPP_NFC_HANDOVER_SEL] command to hostapd
//
// own is DPP record id and uri is NFC tag data uri
func (c *Client) DPPNFCHandoverSel(ctx context.Context, own int, uri string) (string, error) {
	return c.Do(ctx, DPP_NFC_HANDOVER_SEL, fmt.Sprintf("own=%d", own), fmt.Sprintf("uri=%s", uri))
}

type DPPBootstrapGenType string

const (
	DPPBootstrapGenTypeQRCode DPPBootstrapGenType = "qrcode"
	DPPBootstrapGenTypePKEX   DPPBootstrapGenType = "pkex"
	DPPBootstrapGenTypeNFCURI DPPBootstrapGenType = "nfc-uri"
)

// DPPBootstrapGen sends [DPP_BOOTSTRAP_GEN] command to hostapd
//
// channel, mac, info, curve, key, supported_curves and host are optional
func (c *Client) DPPBootstrapGen(
	ctx context.Context,
	typ DPPBootstrapGenType,
	channel string,
	mac string,
	info string,
	curve string,
	key string,
	supported_curves string,
	host string,
) (string, error) {
	args := []string{}
	args = append(args, string(typ))

	if channel != "" {
		args = append(args, fmt.Sprintf("channel=%s", channel))
	}

	if mac != "" {
		args = append(args, fmt.Sprintf("mac=%s", mac))
	}

	if info != "" {
		args = append(args, fmt.Sprintf("info=%s", info))
	}

	if curve != "" {
		args = append(args, fmt.Sprintf("curve=%s", curve))
	}

	if key != "" {
		args = append(args, fmt.Sprintf("key=%s", key))
	}

	if supported_curves != "" {
		args = append(args, fmt.Sprintf("supported_curves=%s", supported_curves))
	}

	if host != "" {
		args = append(args, fmt.Sprintf("host=%s", host))
	}

	return c.Do(ctx, DPP_BOOTSTRAP_GEN, args...)
}

// DPPBootstrapRemove sends [DPP_BOOTSTRAP_REMOVE] command to hostapd
//
// id is DPP record id or *
func (c *Client) DPPBootstrapRemove(ctx context.Context, id string) (string, error) {
	return c.Do(ctx, DPP_BOOTSTRAP_REMOVE, id)
}

// DPPBootstrapGetURI sends [DPP_BOOTSTRAP_GET_URI] command to hostapd
func (c *Client) DPPBootstrapGetURI(ctx context.Context, id int) (string, error) {
	return c.Do(ctx, DPP_BOOTSTRAP_GET_URI, strconv.Itoa(id))
}

// DPPBootstrapInfo sends [DPP_BOOTSTRAP_INFO] command to hostapd
func (c *Client) DPPBootstrapInfo(ctx context.Context, id int) (string, error) {
	return c.Do(ctx, DPP_BOOTSTRAP_INFO, strconv.Itoa(id))
}

// DPPBootstrapSet sends [DPP_BOOTSTRAP_SET] command to hostapd
func (c *Client) DPPBootstrapSet(ctx context.Context, id int, param string) (string, error) {
	return c.Do(ctx, DPP_BOOTSTRAP_SET, strconv.Itoa(id), param)
}

// DPPAuthInit sends [DPP_AUTH_INIT] command to hostapd
//
// own, role, neg_freq, tcp_addr and tcp_port are optional,
// role must be either `configurator` or `enrollee`,
// tcp_addr must be either an ip address or string `from-uri`
func (c *Client) DPPAuthInit(
	ctx context.Context,
	peer int,
	own int,
	role string,
	neg_freq int,
	tcp_addr string,
	tcp_port int,
) (string, error) {
	args := []string{}
	args = append(args, fmt.Sprintf("peer=%d", peer))

	if own > 0 {
		args = append(args, fmt.Sprintf("own=%d", own))
	}

	if role != "" {
		args = append(args, fmt.Sprintf("role=%s", role))
	}

	if neg_freq > 0 {
		args = append(args, fmt.Sprintf("neg_freq=%d", neg_freq))
	}

	if tcp_addr != "" {
		args = append(args, fmt.Sprintf("tcp_addr=%s", tcp_addr))
	}

	if tcp_port > 0 {
		args = append(args, fmt.Sprintf("tcp_port=%d", tcp_port))
	}

	return c.Do(ctx, DPP_AUTH_INIT, args...)
}

// DPPListen sends [DPP_LISTEN] command to hostapd
//
// role is optional and must be either `configurator` or `enrollee`,
// qr_mutual is optional
func (c *Client) DPPListen(ctx context.Context, freq int, role string, qr_mutual bool) (string, error) {
	args := []string{}
	args = append(args, strconv.Itoa(freq))

	if role != "" {
		args = append(args, fmt.Sprintf("role=%s", role))
	}

	if qr_mutual {
		args = append(args, "qr=mutual")
	}

	return c.Do(ctx, DPP_LISTEN, args...)
}

// DPPStopListen sends [DPP_STOP_LISTEN] command to hostapd
func (c *Client) DPPStopListen(ctx context.Context) (string, error) {
	return c.Do(ctx, DPP_STOP_LISTEN)
}

// DPPConfiguratorAdd sends [DPP_CONFIGURATOR_ADD] command to hostapd
//
// net_access_key_curve, curve, key and ppkey are optional
func (c *Client) DPPConfiguratorAdd(
	ctx context.Context,
	net_access_key_curve string,
	curve string,
	key string,
	ppkey string,
) (string, error) {
	args := []string{}

	if net_access_key_curve != "" {
		args = append(args, fmt.Sprintf("net_access_key_curve=%s", net_access_key_curve))
	}

	if curve != "" {
		args = append(args, fmt.Sprintf("curve=%s", curve))
	}

	if key != "" {
		args = append(args, fmt.Sprintf("key=%s", key))
	}

	if ppkey != "" {
		args = append(args, fmt.Sprintf("ppkey=%s", ppkey))
	}

	return c.Do(ctx, DPP_CONFIGURATOR_ADD, args...)
}

// DPPConfiguratorSet sends [DPP_CONFIGURATOR_SET] command to hostapd
//
// net_access_key_curve is optional
func (c *Client) DPPConfiguratorSet(ctx context.Context, id int, net_access_key_curve string) (string, error) {
	args := []string{}
	args = append(args, strconv.Itoa(id))

	if net_access_key_curve != "" {
		args = append(args, fmt.Sprintf("net_access_key_curve=%s", net_access_key_curve))
	}

	return c.Do(ctx, DPP_CONFIGURATOR_SET, args...)
}

// DPPConfiguratorRemove sends [DPP_CONFIGURATOR_REMOVE] command to hostapd
//
// id is DPP record id or *
func (c *Client) DPPConfiguratorRemove(ctx context.Context, id string) (string, error) {
	return c.Do(ctx, DPP_CONFIGURATOR_REMOVE, id)
}

// DPPConfiguratorSign sends [DPP_CONFIGURATOR_SIGN] command to hostapd
//
// conn_status and akm_use_selector are optional,
// conn_status and akm_use_selector must be either 1 or 0
func (c *Client) DPPConfiguratorSign(
	ctx context.Context,
	curve string,
	conf_query bool,
	configurator int,
	conn_status int,
	akm_use_selector int,
) (string, error) {
	args := []string{}
	args = append(args, fmt.Sprintf("curve=%s", curve))

	if conf_query {
		args = append(args, "conf_query=1")
	}

	args = append(args, fmt.Sprintf("configurator=%d", configurator))

	if conn_status > 0 {
		args = append(args, fmt.Sprintf("conn_status=%d", conn_status))
	}

	if akm_use_selector > 0 {
		args = append(args, fmt.Sprintf("akm_use_selector=%d", akm_use_selector))
	}

	return c.Do(ctx, DPP_CONFIGURATOR_SIGN, args...)
}

// DPPConfiguratorGetKey sends [DPP_CONFIGURATOR_GET_KEY] command to hostapd
func (c *Client) DPPConfiguratorGetKey(ctx context.Context, id int) (string, error) {
	return c.Do(ctx, DPP_CONFIGURATOR_GET_KEY, strconv.Itoa(id))
}

// DPPPkexAdd sends [DPP_PKEX_ADD] command to hostapd
//
// tcp_port, tcp_addr, identifier, code, ver and init are optional
func (c *Client) DPPPkexAdd(
	ctx context.Context,
	own int,
	tcp_port int,
	tcp_addr string,
	identifier string,
	code string,
	ver int,
	init bool,
) (string, error) {
	args := []string{}
	args = append(args, fmt.Sprintf("own=%d", own))

	if tcp_port > 0 {
		args = append(args, fmt.Sprintf("tcp_port=%d", tcp_port))
	}

	if tcp_addr != "" {
		args = append(args, fmt.Sprintf("tcp_addr=%s", tcp_addr))
	}

	if identifier != "" {
		args = append(args, fmt.Sprintf("identifier=%s", identifier))
	}

	if code != "" {
		args = append(args, fmt.Sprintf("code=%s", code))
	}

	args = append(args, fmt.Sprintf("ver=%d", ver))

	if init {
		args = append(args, "init=1")
	}

	return c.Do(ctx, DPP_PKEX_ADD, args...)
}

// DPPPkexRemove sends [DPP_PKEX_REMOVE] command to hostapd
//
// id is DPP record id or *
func (c *Client) DPPPkexRemove(ctx context.Context, id string) (string, error) {
	return c.Do(ctx, DPP_PKEX_REMOVE, id)
}

// DPPControllerStart sends [DPP_CONTROLLER_START] command to hostapd
//
// tcp_port, role and qr_mutual are optional,
// role must be either `configurator` or `enrollee`,
func (c *Client) DPPControllerStart(ctx context.Context, tcp_port int, role string, qr_mutual bool) (string, error) {
	args := []string{}

	if tcp_port > 0 {
		args = append(args, fmt.Sprintf("tcp_port=%d", tcp_port))
	}

	if role != "" {
		args = append(args, fmt.Sprintf("role=%s", role))
	}

	if qr_mutual {
		args = append(args, "qr=mutual")
	}

	return c.Do(ctx, DPP_CONTROLLER_START, args...)
}

// DPPControllerStop sends [DPP_CONTROLLER_STOP] command to hostapd
func (c *Client) DPPControllerStop(ctx context.Context) (string, error) {
	return c.Do(ctx, DPP_CONTROLLER_STOP)
}

// DPPChirp sends [DPP_CHIRP] command to hostapd
//
// iter and listen are optional
func (c *Client) DPPChirp(ctx context.Context, own, iter, listen int) (string, error) {
	args := []string{}
	args = append(args, fmt.Sprintf("own=%d", own))

	if iter > 0 {
		args = append(args, fmt.Sprintf("iter=%d", iter))
	}

	if listen > 0 {
		args = append(args, fmt.Sprintf("listen=%d", listen))
	}

	return c.Do(ctx, DPP_CHIRP, args...)
}

// DPPStopChirp sends [DPP_STOP_CHIRP] command to hostapd
func (c *Client) DPPStopChirp(ctx context.Context) (string, error) {
	return c.Do(ctx, DPP_STOP_CHIRP)
}

// DPPRelayAddController sends [DPP_RELAY_ADD_CONTROLLER] command to hostapd
func (c *Client) DPPRelayAddController(ctx context.Context, ip, pkhash string) (string, error) {
	return c.Do(ctx, DPP_RELAY_ADD_CONTROLLER, ip, pkhash)
}

// DPPRelayRemoveController sends [DPP_RELAY_REMOVE_CONTROLLER] command to hostapd
func (c *Client) DPPRelayRemoveController(ctx context.Context, ip string) (string, error) {
	return c.Do(ctx, DPP_RELAY_REMOVE_CONTROLLER, ip)
}

// NANPublish sends [NAN_PUBLISH] command to hostapd
func (c *Client) NANPublish(
	ctx context.Context,
	service_name string,
	ttl time.Duration,
	srv_proto_type int,
	ssi string,
	solicited bool,
	unsolicited bool,
	fsd bool,
) (string, error) {
	args := []string{}
	args = append(args, fmt.Sprintf("service_name=%s", service_name))
	args = append(args, fmt.Sprintf("ttl=%d", int(ttl.Seconds())))
	args = append(args, fmt.Sprintf("srv_proto_type=%d", srv_proto_type))
	args = append(args, fmt.Sprintf("ssi=%s", ssi))

	if solicited {
		args = append(args, "solicited=0")
	}

	if unsolicited {
		args = append(args, "unsolicited=0")
	}

	if fsd {
		args = append(args, "fsd=0")
	}

	return c.Do(ctx, NAN_PUBLISH, args...)
}

// NANCancelPublish sends [NAN_CANCEL_PUBLISH] command to hostapd
func (c *Client) NANCancelPublish(ctx context.Context, publish_id int) (string, error) {
	return c.Do(ctx, NAN_CANCEL_PUBLISH, fmt.Sprintf("publish_id=%d", publish_id))
}

// NANUpdatePublish sends [NAN_UPDATE_PUBLISH] command to hostapd
func (c *Client) NANUpdatePublish(ctx context.Context, publish_id int, ssi string) (string, error) {
	return c.Do(ctx, NAN_UPDATE_PUBLISH, fmt.Sprintf("publish_id=%d", publish_id), fmt.Sprintf("ssi=%s", ssi))
}

// NANSubscribe sends [NAN_SUBSCRIBE] command to hostapd
func (c *Client) NANSubscribe(
	ctx context.Context,
	service_name string,
	active bool,
	ttl,
	srv_proto_type int,
	ssi string,
) (string, error) {
	args := []string{}
	args = append(args, fmt.Sprintf("service_name=%s", service_name))

	if active {
		args = append(args, "active=1")
	}

	args = append(args, fmt.Sprintf("ttl=%d", ttl))
	args = append(args, fmt.Sprintf("srv_proto_type=%d", srv_proto_type))
	args = append(args, fmt.Sprintf("ssi=%s", ssi))

	return c.Do(ctx, NAN_SUBSCRIBE, args...)
}

// NANCancelSubscribe sends [NAN_CANCEL_SUBSCRIBE] command to hostapd
func (c *Client) NANCancelSubscribe(ctx context.Context, subscribe_id int) (string, error) {
	return c.Do(ctx, NAN_CANCEL_SUBSCRIBE, fmt.Sprintf("subscribe_id=%d", subscribe_id))
}

// NANTransmit sends [NAN_TRANSMIT] command to hostapd
func (c *Client) NANTransmit(
	ctx context.Context,
	handle,
	req_instance_id int,
	address net.HardwareAddr,
	ssi string,
) (string, error) {
	args := []string{}
	args = append(args, fmt.Sprintf("handle=%d", handle))
	args = append(args, fmt.Sprintf("req_instance_id=%d", req_instance_id))
	args = append(args, fmt.Sprintf("address=%s", address.String()))
	args = append(args, fmt.Sprintf("ssi=%s", ssi))

	return c.Do(ctx, NAN_TRANSMIT, args...)
}

// GetCapability sends [GET_CAPABILITY] command to hostapd
func (c *Client) GetCapability(ctx context.Context, field string) (string, error) {
	return c.Do(ctx, GET_CAPABILITY, field)
}

// PTKSACacheList sends [PTKSA_CACHE_LIST] command to hostapd
func (c *Client) PTKSACacheList(ctx context.Context) (string, error) {
	return c.Do(ctx, PTKSA_CACHE_LIST)
}

// Driver sends [DRIVER] command to hostapd
func (c *Client) Driver(ctx context.Context, cmd string) (string, error) {
	return c.Do(ctx, DRIVER, cmd)
}

// EnableMLD sends [ENABLE_MLD] command to hostapd
func (c *Client) EnableMLD(ctx context.Context) (string, error) {
	return c.Do(ctx, ENABLE_MLD)
}

// DisableMLD sends [DISABLE_MLD] command to hostapd
func (c *Client) DisableMLD(ctx context.Context) (string, error) {
	return c.Do(ctx, DISABLE_MLD)
}
