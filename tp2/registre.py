"""Registre des diplômes délivrés (SQLite) : code -> nom, moyenne, date, signature, révoqué.

usage :
  python3 registre.py liste
  python3 registre.py revoquer <code>
"""
import hashlib
import sqlite3
import sys
from contextlib import closing

DB = "registre.db"


def _connexion():
    co = sqlite3.connect(DB)
    co.row_factory = sqlite3.Row
    co.execute("CREATE TABLE IF NOT EXISTS diplomes (code TEXT PRIMARY KEY, nom TEXT, "
               "moyenne TEXT, date TEXT, signature TEXT, revoque INTEGER DEFAULT 0)")
    return co


def code_diplome(sig):
    """Code court imprimé sur le diplôme : début de l'empreinte SHA-256 de la signature."""
    return hashlib.sha256(sig).hexdigest()[:16]


def enregistrer(code, nom, moyenne, date, sig_b64):
    with closing(_connexion()) as co, co:                # co : commit automatique
        co.execute("INSERT INTO diplomes (code, nom, moyenne, date, signature) VALUES (?,?,?,?,?)",
                   (code, nom, str(moyenne), date, sig_b64))


def chercher(code):
    with closing(_connexion()) as co:
        ligne = co.execute("SELECT * FROM diplomes WHERE code = ?", (code,)).fetchone()
    return dict(ligne) if ligne else None


def revoquer(code):
    with closing(_connexion()) as co, co:
        return co.execute("UPDATE diplomes SET revoque = 1 WHERE code = ?", (code,)).rowcount == 1


if __name__ == "__main__":
    if sys.argv[1:2] == ["liste"]:
        with closing(_connexion()) as co:
            for l in co.execute("SELECT code, nom, moyenne, date, revoque FROM diplomes"):
                print(tuple(l))
    elif sys.argv[1:2] == ["revoquer"] and len(sys.argv) == 3:
        print("révoqué" if revoquer(sys.argv[2]) else "code inconnu")
    else:
        print(__doc__)