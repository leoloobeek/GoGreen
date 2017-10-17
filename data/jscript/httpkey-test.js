function binToHex(binary) {
    var dom = new ActiveXObject("Microsoft.XMLDOM");
    var el = dom.createElement("tmp");
    el.dataType = "bin.hex";
    el.nodeTypedValue = binary;
    el.removeAttribute("dt:dt");
    return el.nodeTypedValue
}
function getSHA512(bytes) {
    var sha512 = new ActiveXObject("System.Security.Cryptography.SHA512Managed");
    var text = new ActiveXObject("System.Text.ASCIIEncoding");
    var result = binToHex(sha512.ComputeHash_2((text.GetBytes_4(bytes))));
    return result
}

var xHttp = new ActiveXObject("MSXML2.XMLHTTP");
xHttp.Open("GET", "~HKURL~", false);
xHttp.setRequestHeader("User-Agent", "~HKUSERAGENT~");
xHttp.Send();
response = xHttp.responseText;
var key = getSHA512(response);
key = key.substring(0,32);

WScript.Echo("Initial HttpKey  : ~HTTPKEY~");
WScript.Echo("HttpKey Obtained : " + key);