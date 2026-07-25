import hmac, hashlib, time, uuid, urllib.request, json

path = "categories"
ts = str(int(time.time()))
nonce = str(uuid.uuid1()).replace("-", "")
method = "GET"
api_key = "C69BAF41DA5ABD1FFEDC6D2FEA56B"
secret = '~d}$Q7$eIni=V)9\\RK/P.RM4;9[7|@/CA}b~OW!3?EV`:<>M7pddUBL5n|0/*Cn'

src = path + ts + nonce + method + api_key
sig = hmac.new(secret.encode('utf-8'), src.lower().encode('utf-8'), hashlib.sha256).hexdigest()

headers = {
    "api-key": api_key, "accept": "application/vnd.picacomic.com.v1+json",
    "app-channel": "3", "time": ts, "app-uuid": "defaultUuid",
    "nonce": nonce, "signature": sig, "app-version": "2.2.1.3.3.4",
    "image-quality": "original", "app-platform": "android",
    "app-build-version": "45", "user-agent": "okhttp/3.8.1",
    "version": "v1.5.4",
}

req = urllib.request.Request("https://picaapi.picacomic.com/categories", headers=headers)
try:
    resp = urllib.request.urlopen(req, timeout=10)
    print(f"Status: {resp.status}")
    print(f"Body: {resp.read().decode()[:500]}")
except urllib.error.HTTPError as e:
    print(f"HTTP Error: {e.code}")
    print(f"Body: {e.read().decode()[:500]}")
except Exception as e:
    print(f"Error: {e}")
