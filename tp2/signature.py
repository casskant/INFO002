"""Signature RSA-PSS / SHA-256 avec PyCryptodome."""
from Crypto.PublicKey import RSA
from Crypto.Signature import pss
from Crypto.Hash import SHA256


def generer_cles(prive="cle_privee.pem", public="cle_publique.pem", passphrase=None):
    cle = RSA.generate(3072)
    with open(prive, "wb") as f:
        f.write(cle.export_key(passphrase=passphrase, pkcs=8,
                               protection="scryptAndAES128-CBC" if passphrase else None))
    with open(public, "wb") as f:
        f.write(cle.publickey().export_key())


def signer(data, chemin_prive, passphrase=None):
    cle = RSA.import_key(open(chemin_prive, "rb").read(), passphrase=passphrase)
    return pss.new(cle).sign(SHA256.new(data))      # signe l'empreinte, pas les données


def verifier(data, signature, chemin_public):
    cle = RSA.import_key(open(chemin_public, "rb").read())
    try:
        pss.new(cle).verify(SHA256.new(data), signature)
        return True
    except (ValueError, TypeError):
        return False

