package mail

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

// Test MimeHeader decoding, not using iconv
func TestMimeHeaderDecode(t *testing.T) {

	/*
		Normally this would fail if not using iconv
		str := MimeHeaderDecode("=?ISO-2022-JP?B?GyRCIVo9dztSOWJAOCVBJWMbKEI=?=")
		if i := strings.Index(str, "【女子高生チャ"); i != 0 {
			t.Error("expecting 【女子高生チャ, got:", str)
		}
	*/

	str := MimeHeaderDecode("=?utf-8?B?55So5oi34oCcRXBpZGVtaW9sb2d5IGluIG51cnNpbmcgYW5kIGg=?=  =?utf-8?B?ZWFsdGggY2FyZSBlQm9vayByZWFkL2F1ZGlvIGlkOm8=?=  =?utf-8?B?cTNqZWVr4oCd5Zyo572R56uZ4oCcU1BZ5Lit5paH5a6Y5pa5572R56uZ4oCd?=  =?utf-8?B?55qE5biQ5Y+36K+m5oOF?=")
	if i := strings.Index(str, "用户“Epidemiology in nursing and health care eBook read/audio id:oq3jeek”在网站“SPY中文官方网站”的帐号详情"); i != 0 {
		t.Error("\nexpecting \n用户“Epidemiology in nursing and h ealth care eBook read/audio id:oq3jeek”在网站“SPY中文官方网站”的帐号详情\n got:\n", str)
	}
	str = MimeHeaderDecode("=?ISO-8859-1?Q?Andr=E9?= Pirard <PIRARD@vm1.ulg.ac.be>")
	if strings.Index(str, "André Pirard") != 0 {
		t.Error("expecting André Pirard, got:", str)
	}
}

// TestMimeHeaderDecodeNone tests strings without any encoded words
func TestMimeHeaderDecodeNone(t *testing.T) {
	// in the best case, there will be nothing to decode
	str := MimeHeaderDecode("Andre Pirard <PIRARD@vm1.ulg.ac.be>")
	if strings.Index(str, "Andre Pirard") != 0 {
		t.Error("expecting Andre Pirard, got:", str)
	}

}

func TestAddressPostmaster(t *testing.T) {
	addr := &Address{User: "postmaster"}
	str := addr.String()
	if str != "postmaster" {
		t.Error("it was not postmaster,", str)
	}
}

func TestAddressNull(t *testing.T) {
	addr := &Address{NullPath: true}
	str := addr.String()
	if str != "" {
		t.Error("it was not empty", str)
	}
}

func TestNewAddress(t *testing.T) {

	addr, err := NewAddress("<hoop>")
	if err == nil {
		t.Error("there should be an error:", err)
	}

	addr, err = NewAddress(`Gogh Fir <tesst@test.com>`)
	if err != nil {
		t.Error("there should be no error:", addr.Host, err)
	}
}

func TestQuotedAddress(t *testing.T) {

	str := `<"  yo-- man wazz'''up? surprise \surprise, this is POSSIBLE@fake.com "@example.com>`
	//str = `<"post\master">`
	addr, err := NewAddress(str)
	if err != nil {
		t.Error("there should be no error:", err)
	}

	str = addr.String()
	// in this case, string should remove the unnecessary escape
	if strings.Contains(str, "\\surprise") {
		t.Error("there should be no \\surprise:", err)
	}

}

func TestAddressWithIP(t *testing.T) {
	str := `<"  yo-- man wazz'''up? surprise \surprise, this is POSSIBLE@fake.com "@[64.233.160.71]>`
	addr, err := NewAddress(str)
	if err != nil {
		t.Error("there should be no error:", err)
	} else if addr.IP == nil {
		t.Error("expecting the address host to be an IP")
	}
}

func TestEnvelope(t *testing.T) {
	e := NewEnvelope("127.0.0.1", 22)

	e.QueuedId = "abc123"
	e.Helo = "helo.example.com"
	e.MailFrom = Address{User: "test", Host: "example.com"}
	e.TLS = true
	e.RemoteIP = "222.111.233.121"
	to := Address{User: "test", Host: "example.com"}
	e.PushRcpt(to)
	if to.String() != "test@example.com" {
		t.Error("to does not equal test@example.com, it was:", to.String())
	}
	e.Data.WriteString("Subject: Test\n\nThis is a test nbnb nbnb hgghgh nnnbnb nbnbnb nbnbn.")

	addHead := "Delivered-To: " + to.String() + "\n"
	addHead += "Received: from " + e.Helo + " (" + e.Helo + "  [" + e.RemoteIP + "])\n"
	e.DeliveryHeader = addHead

	r := e.NewReader()

	data, _ := ioutil.ReadAll(r)
	if len(data) != e.Len() {
		t.Error("e.Len() is incorrect, it shown ", e.Len(), " but we wanted ", len(data))
	}
	if err := e.ParseHeaders(); err != nil && err != io.EOF {
		t.Error("cannot parse headers:", err)
		return
	}
	if e.Subject != "Test" {
		t.Error("Subject expecting: Test, got:", e.Subject)
	}

}

func TestEncodedWordAhead(t *testing.T) {
	str := "=?ISO-8859-1?Q?Andr=E9?= Pirard <PIRARD@vm1.ulg.ac.be>"
	if hasEncodedWordAhead(str, 24) != -1 {
		t.Error("expecting no encoded word ahead")
	}

	str = "=?ISO-8859-1?Q?Andr=E9?= ="
	if hasEncodedWordAhead(str, 24) != -1 {
		t.Error("expecting no encoded word ahead")
	}

	str = "=?ISO-8859-1?Q?Andr=E9?= =?ISO-8859-1?Q?Andr=E9?="
	if hasEncodedWordAhead(str, 24) == -1 {
		t.Error("expecting an encoded word ahead")
	}

}

func TestParseSubject(t *testing.T) {
	subject := "Bug#1024883: Reconsider versioned dependencies of libpinyin-data"
	e := NewEnvelope("127.0.0.1", 22)

	e.QueuedId = "abc123"
	e.Helo = "helo.example.com"
	e.MailFrom = Address{User: "test", Host: "example.com"}
	e.TLS = true
	e.RemoteIP = "222.111.233.121"
	to := Address{User: "test", Host: "example.com"}
	e.PushRcpt(to)
	if to.String() != "test@example.com" {
		t.Error("to does not equal test@example.com, it was:", to.String())
	}
	e.Data.WriteString(`Delivered-To: kiq1q69q0x@mail.usemail.dev
Received: from 2001:41b8:202:deb:216:36ff:fe40:4002 ([2001:41b8:202:deb:216:36ff:fe40:4002])
	by ztmail.net with ESMTPS id 514cf5f6a4e9abad02c1be32c84ed83b@ztmail.net;
	Sun, 27 Nov 2022 19:03:10 +0000
Received: from localhost (localhost [127.0.0.1])
	by bendel.debian.org (Postfix) with QMQP
	id 9659B20C57; Sun, 27 Nov 2022 19:03:09 +0000 (UTC)
X-Mailbox-Line: From debian-bugs-dist-request@lists.debian.org  Sun Nov 27 19:03:09 2022
Old-Return-Path: <debbugs@buxtehude.debian.org>
X-Original-To: lists-debian-bugs-dist@bendel.debian.org
Delivered-To: lists-debian-bugs-dist@bendel.debian.org
Received: from localhost (localhost [127.0.0.1])
	by bendel.debian.org (Postfix) with ESMTP id 2108020C4A
	for <lists-debian-bugs-dist@bendel.debian.org>; Sun, 27 Nov 2022 19:03:09 +0000 (UTC)
X-Virus-Scanned: at lists.debian.org with policy bank bug
X-Spam-Flag: NO
X-Spam-Score: -4.2
X-Spam-Level:
X-Spam-Status: No, score=-4.2 tagged_above=-10000 required=5.3
	tests=[BAYES_00=-1.9, RCVD_IN_DNSWL_MED=-2.3]
	autolearn=ham autolearn_force=no
Received: from bendel.debian.org ([127.0.0.1])
	by localhost (lists.debian.org [127.0.0.1]) (amavisd-new, port 2525)
	with ESMTP id D0BE3OLHb0gl
	for <lists-debian-bugs-dist@bendel.debian.org>;
	Sun, 27 Nov 2022 19:03:05 +0000 (UTC)
Received: from buxtehude.debian.org (buxtehude.debian.org [IPv6:2607:f8f0:614:1::1274:39])
	(using TLSv1.3 with cipher TLS_AES_256_GCM_SHA384 (256/256 bits)
		key-exchange ECDHE (P-256) server-signature RSA-PSS (2048 bits) server-digest SHA256
		client-signature RSA-PSS (2048 bits) client-digest SHA256)
	(Client CN "buxtehude.debian.org", Issuer "Debian SMTP CA" (not verified))
	by bendel.debian.org (Postfix) with ESMTPS id 4664D20BF0;
	Sun, 27 Nov 2022 19:03:05 +0000 (UTC)
Received: from debbugs by buxtehude.debian.org with local (Exim 4.94.2)
	(envelope-from <debbugs@buxtehude.debian.org>)
	id 1ozMvd-003jpF-Jx; Sun, 27 Nov 2022 19:03:01 +0000
X-Loop: owner@bugs.debian.org
Subject: Bug#1024883: Reconsider versioned dependencies of libpinyin-data
Reply-To: Boyuan Yang <byang@debian.org>, 1024883@bugs.debian.org
Resent-From: Boyuan Yang <byang@debian.org>
Resent-To: debian-bugs-dist@lists.debian.org
Resent-CC: gunnarhj@debian.org, Debian Input Method Team <debian-input-method@lists.debian.org>
X-Loop: owner@bugs.debian.org
Resent-Date: Sun, 27 Nov 2022 19:03:00 +0000
Resent-Message-ID: <handler.1024883.B1024883.1669575614890464@bugs.debian.org>
X-Debian-PR-Message: followup 1024883
X-Debian-PR-Package: src:libpinyin
X-Debian-PR-Keywords:
References: <c0f50339-dd6f-7fbf-126c-4cc50f2fd990@debian.org> <c0f50339-dd6f-7fbf-126c-4cc50f2fd990@debian.org>
X-Debian-PR-Source: libpinyin
Received: via spool by 1024883-submit@bugs.debian.org id=B1024883.1669575614890464
			(code B ref 1024883); Sun, 27 Nov 2022 19:03:00 +0000
Received: (at 1024883) by bugs.debian.org; 27 Nov 2022 19:00:14 +0000
X-Spam-Bayes: score:0.0000 Tokens: new, 16; hammy, 150; neutral, 78; spammy,
	0. spammytokens: hammytokens:0.000-+--H*u:Evolution, 0.000-+--ramacher,
		0.000-+--Ramacher, 0.000-+--soname, 0.000-+--H*ct:application
Received: from mail-qk1-f170.google.com ([209.85.222.170]:34316)
	by buxtehude.debian.org with esmtps (TLS1.3:ECDHE_X25519__RSA_PSS_RSAE_SHA256__AES_128_GCM:128)
	(Exim 4.94.2)
	(envelope-from <073plan@gmail.com>)
	id 1ozMsv-003je0-7K
	for 1024883@bugs.debian.org; Sun, 27 Nov 2022 19:00:14 +0000
Received: by mail-qk1-f170.google.com with SMTP id c2so5891134qko.1
		for <1024883@bugs.debian.org>; Sun, 27 Nov 2022 11:00:12 -0800 (PST)
X-Google-DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
		d=1e100.net; s=20210112;
		h=mime-version:user-agent:organization:references:in-reply-to:date:to
			:from:subject:message-id:x-gm-message-state:from:to:cc:subject:date
			:message-id:reply-to;
		bh=4Ix61NlKLKoTlTTxQv5noTQKK050fyUN4ZBYhXH5DTc=;
		b=I594iYUtjXUdqc6iFzxEH/UsXRIGxCe3vyi8P8v4gQ30DSDqkd6YJJc/61M1/whJ+n
			6pAZfggYyT+d5kcSdoRXpz04PbwmW4ZoMGnOqsABMdxT0XNpixUXtf18FVWPIeq1bQDA
			4u168/ccDWwWJZ2oH16NrJtmPOJjWLkyqFSyDmZvcblGVzS3pMPtVkhY1iMNY68XkPEY
			L9ZoVS5NfLs8qOvxjifrQxpyGG91PAnSlnm3J9Nc/CJO9V7ElE9BriLHefpBdPynJPkd
			XHzW8wOBID6rNruaHgRgAgnqvQg/nl4STtGXZr2YA05mALQbm5TE++QLKBVsh+Ri8cG2
			6Zeg==
X-Gm-Message-State: ANoB5pmQQDRWBFB7IV6nhF5iOG5FD/Zjb4heRMhPNk1LGpmAcfEdvz/1
	KT1qcszEftZaNAx/jhJaGBqiv3Yz03w=
X-Google-Smtp-Source: AA0mqf7IOu/PZySvUIIAWzrxjBORUzyvpJk/oEc3qO2SbPinHsy5tUekcH63thcuJaV8dT75+ZiBsw==
X-Received: by 2002:a37:64c8:0:b0:6fa:182:f2f1 with SMTP id y191-20020a3764c8000000b006fa0182f2f1mr42577708qkb.710.1669575609674;
		Sun, 27 Nov 2022 11:00:09 -0800 (PST)
Received: from gaolab004.ece.pitt.edu ([136.142.25.52])
		by smtp.gmail.com with ESMTPSA id h11-20020ac8714b000000b00342f8d4d0basm5681119qtp.43.2022.11.27.11.00.08
		for <1024883@bugs.debian.org>
		(version=TLS1_3 cipher=TLS_AES_256_GCM_SHA384 bits=256/256);
		Sun, 27 Nov 2022 11:00:08 -0800 (PST)
Message-ID: <23b0fa837b039c3af5f8faafbb07f96e04437727.camel@debian.org>
From: Boyuan Yang <byang@debian.org>
To: 1024883@bugs.debian.org
Date: Sun, 27 Nov 2022 14:00:07 -0500
In-Reply-To: <c0f50339-dd6f-7fbf-126c-4cc50f2fd990@debian.org>
Organization: Debian Project
Content-Type: multipart/signed; micalg="pgp-sha512";
	protocol="application/pgp-signature"; boundary="=-6+C7LQxijPPbb0XpHyN5"
User-Agent: Evolution 3.46.1-1
MIME-Version: 1.0
X-Debian-Message: from BTS
X-Mailing-List: <debian-bugs-dist@lists.debian.org> archive/latest/1745699
X-Loop: debian-bugs-dist@lists.debian.org
List-Id: <debian-bugs-dist.lists.debian.org>
List-URL: <https://lists.debian.org/debian-bugs-dist/>
List-Post: <mailto:debian-bugs-dist@lists.debian.org>
List-Help: <mailto:debian-bugs-dist-request@lists.debian.org?subject=help>
List-Subscribe: <mailto:debian-bugs-dist-request@lists.debian.org?subject=subscribe>
List-Unsubscribe: <mailto:debian-bugs-dist-request@lists.debian.org?subject=unsubscribe>
Precedence: list
Resent-Sender: debian-bugs-dist-request@lists.debian.org


--=-6+C7LQxijPPbb0XpHyN5
Content-Type: text/plain; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

X-Debbugs-CC: gunnarhj@debian.org

Hi,

=E5=9C=A8 2022-11-27=E6=98=9F=E6=9C=9F=E6=97=A5=E7=9A=84 13:31 +0100=EF=BC=
=8CGunnar Hjalmarsson=E5=86=99=E9=81=93=EF=BC=9A
> Package: src:libpinyin
> Version: 2.7.92-2
>=20
> Hi all!
>=20
> Upstream made a SONAME bumb:
>=20
> https://github.com/libpinyin/libpinyin/commit/2f52299e
>=20
> While I don't really understand the reason for it, the resulting=20
> paperwork in Debian is now done.
>=20
> At the transition bug we had this conversation:
>=20
> On 2022-11-26 15:44, Sebastian Ramacher wrote:
> > On 2022-11-25 14:52:12 +0100, Gunnar Hjalmarsson wrote:
> > > I notice that libpinyin has not yet migrated, even though the 2
> > > days delay is over. Is that because Britney waits for the
> > > dependencies to be migration ready too, or is it because this bug
> > > is not closed yet?
> >=20
> > It has not migrated yet because the shared library packages have
> > strictly versioned dependency on libpinyin-data. Hence, migrating
> > libpinyin to testing would currently render some packages
> > uninstallable in testing.
> >=20
> > Ideally, this dependency would be relaxed if possible so that this
> > won't be an issue for the next libpinyin transition. For this one,
> > all the reverse dependencies and libpinyin need to migrate together.
>=20
> I take it that Sebastian talks about:
>=20
> Depends: libpinyin-data (=3D ${binary:Version})
>=20
> in libpinyin15 and libzhuyin15.
>=20
> Question: Would it be an option to change those dependencies to not be=
=20
> versioned in order to make future transitions easier, or are there=20
> reasons for keeping it as it is?

It really depends on whether mismatched libpinyin library + libpinyin-data
will cause big troubles. If mismatched versions will cause crash, the
relaxed dependency should make no sense. Current implementation is the most
conservative one, yet it won't cause obvious troubles.

Thanks,
Boyuan Yang

--=-6+C7LQxijPPbb0XpHyN5
Content-Type: application/pgp-signature; name="signature.asc"
Content-Description: This is a digitally signed message part

-----BEGIN PGP SIGNATURE-----

iQIzBAABCgAdFiEEfncpR22H1vEdkazLwpPntGGCWs4FAmODs7cACgkQwpPntGGC
Ws5HIA//Uj+Sp3m1IzL4gMtPUJxksNG9TpG0qDR8k1pZdydwZeszJrxCn+5S9H5F
InHIhZICFkW+5AwebhV/aUHrEFbBO0fOWYwQf+HEAYxVcOl3EjPXE5p7EAQ+V5eQ
KcDy1rmpr1b7FppPDIIrVO70eOJScBbMS5FctAExOupTDcvKAYMYlqfCcFpv6goY
hnseQ4854fUiPgayVOYIIkvQyCA8rPMtL4dg0StRgOUFr3urdCM4k3bG8UMAZ/By
hFHRT+l6HfNEs8M/Yd8vJd0N6YP/nO4OK3/hvMuyCyqc8ZgfQDDUgJuwLGAPTiJn
K+b+DApf1Xc0NOfWIRxNvyUbCWXTaKJilPA96iwfGoFFUrRpP2Ogmvt7XXk0NQd2
TbeW+DnrphVAI0TiiERq56DIoHq100IIVD+t/ZS+pTPKF57DBKoDiVQjFw//Y0np
nGDka4IErbefLzEZZgf/UBqyLEGgYa6/RvLTH5sLPhnEJwBIRQF2++Xf3LSgTZab
PMADkFwSxCaT8KL+vjCM+XJw5vqUYbwiFR5awxSLpVKFOSvEy1KZO8cXIkv86mgd
tJu+vocoM2OQkfsj2Qykl3CitUo0bGel86MtOScvYmm83bbxLlrYOjPHZaA5NgjD
LsWYuyOmhCC/mxdsaxZisUcMbnZLZ29ovoaTVL1gHw+K6OG1R8c=
=TFff
-----END PGP SIGNATURE-----

--=-6+C7LQxijPPbb0XpHyN5--

	`)

	addHead := "Delivered-To: " + to.String() + "\n"
	addHead += "Received: from " + e.Helo + " (" + e.Helo + "  [" + e.RemoteIP + "])\n"
	e.DeliveryHeader = addHead

	r := e.NewReader()

	data, _ := ioutil.ReadAll(r)
	if len(data) != e.Len() {
		t.Error("e.Len() is incorrect, it shown ", e.Len(), " but we wanted ", len(data))
	}
	if err := e.ParseHeaders(); err != nil && err != io.EOF {
		t.Error("cannot parse headers:", err)
		return
	}
	if e.Subject != subject {
		t.Errorf("Subject expecting: %s, got: %s", subject, e.Subject)
	}
}
