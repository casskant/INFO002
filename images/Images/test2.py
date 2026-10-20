from PIL import Image
from PIL import ImageFont
from PIL import ImageDraw


def add_text(img, text, x, y):
    draw = ImageDraw.Draw(img)                  # objet "dessin" dans l'image
    font = ImageFont.truetype("sans.ttf", 32)   # police à utiliser
    draw.text((x, y), text, "white", font)       # ajout du texte


def main(filename, text, output):
    img = Image.open(filename)  # ouverture de l'image contenue dans un fichier
    add_text(img, text, 110, 110)
    img.save(output)            # sauvegarde de l'image obtenue dans un autre fichier


if __name__ == "__main__":
    import sys
    if len(sys.argv) != 4:
        print("usage: {} image msg output".format(sys.argv[0]))
        sys.exit(1)
    main(sys.argv[1], sys.argv[2], sys.argv[3])
