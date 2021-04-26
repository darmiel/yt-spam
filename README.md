# yt-spam
Small tool that tries to detect bots in comments.

## Checks
There are currently the following checks:
* [🔁 Format-Spam](#format-spam)
* [🐈 Copy-Cat](#copy-cat)
* ✍️ Blacklist
    * [✍️ Name-Blacklist](#name-blacklist)
    * [✍️ Body-Blacklist](#comment-blacklist)

### Format-Spam
Searches for recurring formatted words; example:
```
UserA) Google this! *ABC1234*
UserB) *ABC1234* <- google that!
UserC) no scam!!!11: *ABC1234*
```
Any comment containing such a word (`*ABC1234*`) will be marked as malicious.

![Screenshot 2021-04-26 at 19 23 20](https://user-images.githubusercontent.com/71837281/116124674-e1fa2500-a6c4-11eb-9d0a-be23dfb3906d.png)


### Copy-Cat
Checks (long) comments for duplicates

https://user-images.githubusercontent.com/71837281/115129156-ba7bcc00-9fe3-11eb-961c-ebef3928906c.mov


```
🐈 COPY-CAT Shakina Eprillia copied Wendy Williams w/ why am i watching this? + UgyXwt3AMowF3cD1B5h4AaABAg , - Ugx3fOsQdPHLQFV9dDJ4AaABAg ]
🐈 COPY-CAT Shyamlee singh copied Ray w/ Please be my math teacher + Ugw9JPlsDQfqZ_lRS4B4AaABAg , - Ugx58npioj9isYLYzdp4AaABAg ]
🐈 COPY-CAT Kitsuné Noir copied xavier mosqueda w/ I don’t get it + UgxSfUwd5HjQPB1t61p4AaABAg , - UgxEB032vNWbSSPCfs94AaABAg ]
🐈 COPY-CAT darkleenk1 copied Jiya w/ what's a tetrahedron + UgwZt5eMiqqaIAU_0DJ4AaABAg , - UgwOKE6_38C4VEhp6L94AaABAg ]
🐈 COPY-CAT you got no jams copied Asouthindiandude w/ *brain has left the chat* + UgzekPSyiQZD8dVisIx4AaABAg , - UgxeUchqqf6G6qhhlzF4AaABAg ]
🐈 COPY-CAT Carbo FLx copied ben w/ why am i watching this + UgxLxPjDeKnGm37vxPJ4AaABAg , - UgyStifNlp2ZJpo7Nst4AaABAg ]
🐈 COPY-CAT King Zig copied Dbeaumier 4010 w/ I like your funny words magic man + Ugymms_4TB4lBSBrg-B4AaABAg , - UgyfCfFor1xLK1Wqns14AaABAg ]
🐈 COPY-CAT Max copied Jonathan Avila w/ Depends on size of center + UgwxAns6hggoo2EcIad4AaABAg , - UgzOVERaBXIEoNu_5w94AaABAg ]
🐈 COPY-CAT AruBoii copied Dizzy w/ Dawg why did the YouTube algorit... + UgyvFa0AdB8_F5zhSbR4AaABAg , - UgyXPpdtOXZe81Vprd54AaABAg ]
```

### Name-Blacklist
Checks names for blacklisted words (`data/input/name-blacklist.txt`)
![Screenshot 2021-04-21 at 16 51 51](https://user-images.githubusercontent.com/71837281/115574446-eb991c80-a2c1-11eb-96cb-8580e306fcf3.png)

### Comment-Blacklist
Checks comment-bodies for blacklisted words (`data/input/comment-blacklist.txt`)
