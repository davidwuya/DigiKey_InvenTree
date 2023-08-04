from blabel import LabelWriter
import logging
import os
import argparse

def write_labels(ManufacturerPartNumber: str, LimitedTaxonomy: str, ProductDescription: str) -> None:
    # find fonttools logging and disable it
    logging.getLogger("fontTools").setLevel(logging.CRITICAL)
    label_writer = LabelWriter("template.html", default_stylesheets=("style.css",))
    # longest allowable string is maxlength characters
    maxlength = 14
    ManufacturerPartNumber = (
        ManufacturerPartNumber[:maxlength]
        if len(ManufacturerPartNumber) > maxlength
        else ManufacturerPartNumber
    )
    Category = (
        LimitedTaxonomy[:maxlength]
        if len(LimitedTaxonomy) > maxlength
        else LimitedTaxonomy
    )
    Description1 = ""
    Description2 = ""
    if len(ProductDescription) > maxlength:
        Description_List = ProductDescription.split(" ")
        # fit as many words as possible into first line
        for word in Description_List:
            if len(Description1) + len(word) > maxlength:
                break
            Description1 += word + " "
        # if there are still words left, fit as many as possible into second line and store in Description2
        for word in Description_List[len(ProductDescription.split(" ")) - 1 :]:
            if len(Description2) + len(word) > maxlength:
                break
            Description2 += word + " "
    records = [
        dict(
            ManufacturerPartNumber=ManufacturerPartNumber,
            Category=Category,
            Description_Line1=Description1,
            Description_Line2=Description2,
        ),
    ]
    fname = f"{ManufacturerPartNumber}.pdf"
    try:
        label_writer.write_labels(records, target=fname)
        logging.info(f"Labels written to {fname}")
        # move to labels folder
        os.rename(fname, os.path.join("labels", fname))
        # open the PDF
        os.startfile(os.path.join("labels", fname))
    except FileExistsError:
        logging.error(f"File {fname} already exists.")
        pass

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Write labels for a product.")
    parser.add_argument(
        "-m",
        "--ManufacturerPartNumber",
        type=str,
        help="Manufacturer Part Number",
        required=True,
    )
    parser.add_argument(
        "-l",
        "--LimitedTaxonomy",
        type=str,
        help="Limited Taxonomy",
        required=True,
    )
    parser.add_argument(
        "-d",
        "--ProductDescription",
        type=str,
        help="Product Description",
        required=True,
    )
    args = parser.parse_args()
    write_labels(args.ManufacturerPartNumber, args.LimitedTaxonomy, args.ProductDescription)