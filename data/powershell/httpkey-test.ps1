function Get-SHA512Hash {
        param($b)
        $sha512 = New-Object System.Security.Cryptography.SHA512CryptoServiceProvider
        return [System.BitConverter]::ToString($sha512.ComputeHash([System.Text.Encoding]::ASCII.GetBytes($b))).ToLower() -replace "-",""
}
[System.Net.ServicePointManager]::Expect100Continue=0;
$wc = New-Object Net.WebClient
$wc.Headers.Add('User-Agent','~HKUSERAGENT')
$wc.Proxy=[System.Net.WebRequest]::DefaultWebProxy
$wc.Proxy.Credentials = [System.Net.CredentialCache]::DefaultNetworkCredentials
$resp = $wc.DownloadString('~HKURL~')   

$key = (Get-SHA512Hash $resp)[0..31] -join ""

Write-Host "Initial HttpKey  : ~HTTPKEY~"
Write-Host "HttpKey Obtained : $key"
