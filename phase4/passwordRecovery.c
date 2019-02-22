
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char* red_herring = "CTF{str1ngs_w0nt_s4v3_u}";
// admin_password = CTF{d4mn_u_sm4rt_gurl}
char* admin_mash = "a8mTCX|}.Y.J|Xs,.1Ys{+";
// ceclia_password = DanielIsTheCoolest
char* ceclia_mash = "4^%cfjLu>Vf_\"\"jfuw";

// generates a random cipher alphabet
// len is the length of the password or password mash. This is the random seed
// cipher is a char[] that will contain the random alphabet from index 32 to 127
void poly(int len, char* cipher) {
    srand(len);
    for (int i=32; i<128; i++) {
        cipher[i] = 0;
    }
    for (int c=32; c<128; c++) {
        int i = (rand() % 96) + 32;
        // slide i up until it finds an empty slot for c
        while (cipher[i] != 0) {
            i = ((i - 31) % 96) + 32;
        }
        cipher[i] = (char)c;
    }
}

// Mashes up a password
char* mash_it_up(const char* password) {
    int len = strlen(password);
    srand(len);
    char cipher[128];
    poly(len, cipher);
    char* mash = (char*)malloc(len * sizeof(char));
    for (int i=0; i<len; i++) {
        mash[i] = cipher[password[i]];
    }
    mash[len] = 0;
    return mash;
}

// Unmashes a password
char* mash_it_down(const char* mash) {
    int len = strlen(mash);
    srand(len);
    char cipher[128];
    poly(len, cipher);
    char* password = (char*)malloc(len * sizeof(char));
    for (int i=0; i<len; i++) {
        for (int j=32; j<128; j++) {
            if (cipher[j] == mash[i]) {
                password[i] = (char)j;
                break;
            }
        }
    }
    password[len] = 0;
    return password;
}

int main(int argc, const char* argv[]) {
    if (argc != 2) {
        printf("Please specify exactly one argument for the password\n");
        return 1;
    }
    char* mashed_password = mash_it_up(argv[1]);
    if (strcmp(mashed_password, admin_mash) != 0) {
        printf("Incorrect password!\n");
        return 1;
    }
    free(mashed_password);
    printf("Admin password accepted!\n");
    char* unmashed_password = mash_it_down(ceclia_mash);
    printf("Ceclia's password: %s\n", unmashed_password);
    return 0;
}


