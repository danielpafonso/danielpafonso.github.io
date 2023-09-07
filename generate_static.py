""" Generate static webpages """
import os
import shutil

from datetime import datetime

def read_qoutes_file(filename):
    """ read qoute files and return list """
    # read file
    with open(filename, "r", encoding="utf8") as ifp:
        raw = ifp.read().strip().split("\n")
    # process data
    qoutes = [line.split("|") for line in raw]
    qoutes.sort(key=lambda x: x[0])
    return qoutes

def prep_static_folder(folder, imgs_folder="img"):
    """ create static folder and copy the images to it """
    # create folder
    print("Creating static folder")
    os.makedirs(folder, exist_ok=True)
    # copy images
    print("Copying images to static folder")
    shutil.copytree(imgs_folder, os.path.join(folder, imgs_folder))

def generate_index(data:list, input_folder, output_folder, filename="index.html"):
    """ Generate static page for today qoute """
    print("Generating index page")
    # read template
    template_file = os.path.join(input_folder, f"template_{filename}")
    with open(template_file, "r", encoding="utf8") as ifp:
        template = ifp.read()

    # set defaut qoute as the last one
    phrase = data[-1][1]
    author = data[-1][2]
    # get current qoute
    current_day = datetime.utcnow().strftime("%Y-%m-%d")
    for x in data:
        if x[0] == current_day:
            phrase = x[1]
            author = x[2]
            break
    
    # replace values
    template = template.format(phrase=phrase, author=author)

    # write static
    static_file = os.path.join(output_folder, filename)
    with open(static_file, "w", encoding="utf8") as ofp:
        ofp.write(template)

def generate_list(data, input_folder, output_folder, filename="list.html"):
    """ Generate static page for list of all qoutes """
    print("Generating list page")
    template_file = os.path.join(input_folder, f"template_{filename}")
    with open(template_file, "r", encoding="utf8") as ifp:
        template = ifp.read().split("\n")

    # get row template
    list_idx = 0
    row_template = ""
    for i, line in enumerate(template):
        if "{qlist}" in line:
            list_idx = i
            row_template = line
            break
    
    # get qoute rows
    rows = [
        str(row_template).format(qlist="", date=qoute[0], phrase=qoute[1], author=qoute[2])
        for qoute in data
        ]
    # replace values
    template[list_idx] = "\n".join(rows)

    # write static
    static_file = os.path.join(output_folder, filename)
    with open(static_file, "w", encoding="utf8") as ofp:
        ofp.write("\n".join(template))

# def generate_duplicates(data, input_folder, output_folder, filename="duplicates.html"):
#     """ Generate static page for duplicates qoutes"""

if __name__ == "__main__":
    filepath = "quotes.dsv"
    template_folder = "templates"
    static_folder = "static"
    
    # read qoute file
    quotes = read_qoutes_file(filepath)

    # prep output folder
    prep_static_folder(static_folder)

    # generate static page from templates
    generate_index(quotes, template_folder, static_folder)
    generate_list(quotes, template_folder, static_folder)
    ## dummy duplicate html
    filename = "duplicates.html"
    shutil.copyfile(filename, os.path.join(static_folder, filename))
    # generate_duplicates(quotes, template_folder, static_folder)
