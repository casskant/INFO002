"""Stéganographie LSB : on cache des octets dans le bit de poids faible du rouge."""
import struct
from PIL import Image


def capacite(img):
    """Nombre d'octets cachables (1 bit par pixel, moins 4 octets pour la taille)."""
    return img.width * img.height // 8 - 4


def _bits(data):
    for octet in data:
        for i in range(7, -1, -1):          # bit de poids fort d'abord
            yield (octet >> i) & 1


def cacher(img, data):
    """Renvoie une copie de img contenant data. Format caché : [taille sur 4 octets][data]."""
    if len(data) > capacite(img):
        raise ValueError(f"message trop long ({len(data)} > {capacite(img)} octets)")
    img = img.convert("RGB")                # copie, sans canal alpha
    px = img.load()
    for k, bit in enumerate(_bits(struct.pack(">I", len(data)) + data)):
        y, x = divmod(k, img.width)         # k-ème pixel, ligne par ligne
        r, g, b = px[x, y]
        px[x, y] = ((r & ~1) | bit, g, b)   # on remplace le LSB du rouge
    return img


def _lire(px, largeur, debut, nb_octets):
    """Lit nb_octets en partant du bit numéro `debut`."""
    out = bytearray()
    for o in range(nb_octets):
        v = 0
        for i in range(8):
            y, x = divmod(debut + o * 8 + i, largeur)
            v = (v << 1) | (px[x, y][0] & 1)
        out.append(v)
    return bytes(out)


def extraire(img, nb_bytes=None):
    """Récupère le message caché (taille lue dans l'en-tête sauf si nb_bytes est donné)."""
    img = img.convert("RGB")
    px = img.load()
    if nb_bytes is None:
        nb_bytes = struct.unpack(">I", _lire(px, img.width, 0, 4))[0]
        if nb_bytes > capacite(img):
            raise ValueError("aucun message valide dans cette image")
    return _lire(px, img.width, 32, nb_bytes)


# noms utilisés par diplome.py (même format : en-tête de taille + message)
cacher_msg = cacher


def extraire_msg(img):
    return extraire(img)