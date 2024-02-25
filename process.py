""" process quotes to revriteve insigths """
import re

def read_file(filename):
    """ read and process data file """

    # read file
    with open(filename, "r", encoding="utf8") as ifp:
        raw = ifp.read().strip().split("\n")

    # process data
    data = [line.split("|") for line in raw]

    return data

def groupby_author(data):
    """ process data to group qoutes by author """
    groupby = {}
    for x in data:
        # extract author
        author = x[2].split(" - ")[0].split(", ")[0]
        if author in groupby:
            groupby[author].append(x[1])
        else:
            groupby[author] = [x[1]]

    for author, quotes in groupby.items():
        if len(quotes) > 1:
            print(author)
            print("", "\n ".join(sorted(quotes)))
            print()

if __name__ == "__main__":
    data = read_file("quotes.dsv")
    groupby_author(data)
