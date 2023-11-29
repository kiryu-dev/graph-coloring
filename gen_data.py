import uuid

filename = str(uuid.uuid4())
file = open(filename, "w")

colors = ['R', 'Y', 'B']
VERTEX_COUNT = 999
buf = [f'{VERTEX_COUNT},{VERTEX_COUNT}']

for i in range(1, VERTEX_COUNT):
    buf.append(f'{i},{i+1}')
buf.append(f'{VERTEX_COUNT},1')

for i in range(VERTEX_COUNT):
    buf.append(f'{i+1},{colors[i % len(colors)]}')

file.write('\n'.join(buf))
file.close()