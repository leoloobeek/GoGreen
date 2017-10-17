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
dim wc, resp, key
set wc = CreateObject("MSXML2.XMLHTTP")
wc.Open "GET", "~HKURL~", false
wc.setRequestHeader "User-Agent", "~HKUSERAGENT~"
wc.Send
resp = wc.responseText
key = Mid(getSHA512(resp), 1, 32)

WScript.Echo "Initial HttpKey  : ~HTTPKEY~"
WScript.Echo "HttpKey Obtained : " & key