int strrindex(char line[], char t) {
  int i = 0;
  int fidx = -1;

  for (i; line[i] != '\0'; i++) {
    if (line[i] == t) fidx = i;
  }

  return fidx;
}
