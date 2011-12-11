#include <stdio.h>

#define MAX_LINE_LEN 1000

int get_line(char line[], int max);
int strrindex(char line[], char t);

int main() {
  char line[MAX_LINE_LEN];
  int found = 0;

  int len = get_line(line, MAX_LINE_LEN);

  printf("Index: %i\n", strrindex(line, getchar()));
}

int get_line(char line[], int max) {
  int c;
  int i = 0;

  while (--max > 0 && (c=getchar()) != EOF && c != '\n')
    line[i++] = c;
  if (c == '\n')
    line[i++] = c;
  line[i] = '\0';

  return i;
}
