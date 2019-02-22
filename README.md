# ceclia-ctf

## Phase 1: Playfair cipher

The message is encrypted with the Playfair cipher with no key. The message will include a hint: `STEGANOGRAPHY`

The ciphertext is handwritten, then photographed. Phase 2 is concealed in the image file. The image file is shared over telegram with the following hint: `play fair ;)`

No flags are included at this stage.

## Phase 2: Steganography

A .zip file is included in the image file from Phase 1. The .zip file contains a README describing the CTF. Basically, the primary objective is figuring out how to submit the flags. The secondary objective is to find all the flags and submit them. Instead of marking progress, flags are dead ends.

The .zip file also contains a directory which contains a number of files which contain a number of ciphertexts. The .zip file also contains a Python program used to generate the ciphertexts. RSA is used to encrypt the same message every time. A different private key is used for every file, and the keys are not given. Hastad's Broadcast Attack must be used to obtain the message. 

The .zip file also contains another .zip file which contains a disclaimer and another .zip file which contains [this .pdf](https://crypto.stanford.edu/~dabo/papers/RSA-survey.pdf). This paper describes a number of RSA vulnerabilities, including Hastad's Broadcast Attack.

The .zip file also contains a file `.flag` that contains the flag. Because it starts with `.` it will not show up on a normal `ls`.

#### Flags

The first flag is contained in the `.flag` file.

## Phase 3: Logging in

The message from Phase 2 is a URL (possibly shortened) to a Raspberry pi hosting my CTF server application (app). The app will automatically redirect to the login page without valid credentials or a hacked auth token. The auth token is a base64 encoded JSON object stored in a cookie:

```
{
    username: string,
    password: string,
    authorized: bool
}
```

The app automatically sets the cookie appropriately when login attempts are made.

The first step is to set authorized=true in the auth token. The app will then allow access to the app pages: Progress, Leaderboard, Admin, and Account.

The Progress page contains all the information about the currently solved flags for the current user. Hints can be requested here.

The Leaderboard page lists all the users and their current completion percentages. There are only two users: ceclia and d4ni3l. d4ni3l is at 100%. This page is where Cecilia should find the correct username for her to use.

The Admin page is only rendered when the user is an admin, which is only true for d4ni3l. The password does not have to be valid for the page to render. The page contains a link to the passwordRecovery executable used in Phase 4. The page also contains a link to reset all user progress, intended for debugging and testing. The password must be valid for the link to work.

The Account page contains a sign-out link. Clicking it will reset the auth token and boot the user back to the login page. 

Subpages are accessed via URL parameters: `/app?page={page}` Each page also has an optional message box which can be used to communicate with the user. Actions are all done through GETs.

#### Flags

The second flag can be found by checking robots.txt on the app.

The third flag can be found by inspecting the JS of the login page.

The fourth flag can be found by specifying an undefined page value in the URL. 

## Phase 4: Reversing

By setting the username to d4ni3l in the auth token, the Admin page is available. The admin page includes a link to an executable file which is advertised as a password recovery tool. The tool is written in C. The tool accepts a single argument, nominally the admin password. It is not actually the password for d4ni3l's app login. If the password is correct, it will print the password to ceclia's app login.

Inside the tool are two functions, `mash_it_up` and `mash_it_down`. These functions encrypt or decrypt a string using a randomly generated substitution cipher. Only characters betwenn 32 and 128 are mapped. The length of the string is used as the seed for the random cipher generation. 

Inside the tool are two mashed up passwords, the admin password and ceclia's password. The validation is done by calling `mash_it_up` on the argument and comparing it to the mashed up admin password. ceclia's password is printed by calling `mash_it_down` on the mashed up ceclia's password. 

#### Flags

The fifth flag can be found by running `strings` on the passwordRecovery tool.

The sixth flag is the admin password. It can be found by calling `mash_it_up` on the mashed up admin password.

## Phase 5: Submitting Flags

Once ceclia's password is obtained, she can log in normally (or set the auth token) and submit any flags discovered so far. The Progress page has the flags sorted so that once a few flags have been found, there is a general idea of where to look for the others. Hints can also be requested as necessary.


