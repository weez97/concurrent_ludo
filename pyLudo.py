import random

players = []
for i in range(4):
  players.append([-1, -1, -1, -1])

map = []
mSize = 57
for i in range(1, mSize):
  if random.random() < 1/5:
    map.append(0)
  else:
    map.append(1)

roll = [0, 0, 0]
won = False
c = 0

while not won:
  c += 1
  print("\nRonda: ", c)
  for i in players:
    roll[0] = random.randint(1, 6)
    roll[1] = random.choice([-1, 1])
    roll[2] = random.randint(1, 6)
    if roll[0] == roll[2] and -1 in i:
      for idx, j in enumerate(i):
        if j == -1:
          i[idx] = 0
          break
    else:
      move = roll[0] + roll[1] * roll[2]
      for idx, j in enumerate(i):
        if j < mSize and j != -1:
          if j + move >= mSize - 1:
            i[idx] = mSize
            break
          if map[j + move]:
            i[idx] = j + move
            break
    print(i)
  for i in players:
    if i == [mSize, mSize, mSize, mSize]:
      won = True
print("\nFinalizado - ", players)