import hmac, hashlib

path = "categories"
ts = "1784960469"
nonce = "5ccc086fc29d4c659059e538feb86be1"
method = "GET"
api_key = "C69BAF41DA5ABD1FFEDC6D2FEA56B"

secret = '~d}$Q7$eIni=V)9\\RK/P.RM4;9[7|@/CA}b~OW!3?EV`:<>M7pddUBL5n|0/*Cn'
print(f"Secret repr: {repr(secret)}")

src = path + ts + nonce + method + api_key
data = src.lower().encode('utf-8')
sig = hmac.new(secret.encode('utf-8'), data, hashlib.sha256).hexdigest()
print(f"Python sig: {sig}")

go_sig = "a56e2cd7cd04f2f077a31b391594ef9fabcde020cf0741c2a264ed131f229bad"
print(f"Match: {sig == go_sig}")
