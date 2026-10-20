from PIL import Image


def invert_line(img):
    m = img.height // 2         # milieu de l'image
    pixels = img.load()         # tableau des pixels

    for x in range(0, img.width):
        r, g, b = pixels[x, m]  # on récupère les composantes RGB du pixel (x,m)
        r = r ^ 0b11111111      # on les inverse bit à bit avec un XOR
        g = g ^ 0b11111111      # ...
        b = b ^ 0b11111111      # ...
        pixels[x, m] = r, g, b  # on remet les pixels inversés dans le tableau


def main(filename, output):
    img = Image.open(filename)  # ouverture de l'image contenue dans un fichier
    invert_line(img)
    img.save(output)            # sauvegarde de l'image obtenue dans un autre fichier


if __name__ == "__main__":
    import sys
    if len(sys.argv) != 3:
        print("usage: {} image output".format(sys.argv[0]))
        sys.exit(1)
    main(sys.argv[1], sys.argv[2])
