<?lua
-- Allow sessions. This needs to be called before anything else.
session_start()

-- Session code. Allow the user to log in if needed!
if REQUEST("action") == "login" then
	session_set("Member")
end
if REQUEST("action") == "logout" then
	session_set("")
end

-- If the user is logged in, show their username!
if session_get("username") == "" then
	print "Welcome Guest! <a href=index.lua?action=login>Login</a><br />"
else
	print("Welcome " .. session_get("username") .. "! <a href=index.lua?action=login>Login</a><br />")
end
print("Serving media files")
?>

<img src="img/dragon.png" /><br />

<?lua
-- Testing MySQL
local userExists = mysql_queryfind("accounts", "username", "sample");
vardump(next(userExists) ~= nil);
mysql_addrow("accounts", {"username" = username, "sha_pass_hash" = sha1(username .. ":" .. password), "email" = email})

?>
