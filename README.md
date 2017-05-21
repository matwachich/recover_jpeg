# Recover JPEG
When you have a JPEG pictures that is corrupted because the begining of the file is lost (corrupted, encrypted...), this tool can help!

It works by extracting the valid data (from 0xFFDA Start Of Scan) from the corrupted file and append it to the valid headers (initial data) of another valid file (model).

You'll have to provide models pictures files taken with the same camera with different resolutions, orientations and quality settings.

# How to use
Create a folder "__models__" asside the executable and fill it with the model files. Name them carefully because the model name is appended to the recovered picture (ex. s7-1080-paysage.jpg).

Then, just drag the corrupted file on the executable. It will try to reconstruct jpeg files with all the models provided. There will be many invalid files asside the corrupted picture, but hopefully there will be one valid file.

# How it works
(I'm absolutly not a specialist, so excuse me if some information is wrong or inaccurate)

The essential part of a JPEG picture is what comes after the 0xFFDA (Start Of Scan) marker. Fortunatly this is the last (and biggest part) or a file, so when the begining of a file is corrupted, the SOS is always here!

What does this tool is extract the SOS data of a corrupted picture, append it the the other parts of a valid file, and tries to parse the result as a new JPEG file.

# Why I did this tool
I'm a victim of a ransomware attack! All my personal pics were gone!

But when I saw that only the first 10kb of each file is encrypted, I was sure there will be a mean to recover my pictures, and indeed there was!

This is why, when no SOS marker is found in the encrypted file, this tool extracts all but the first 10kb and tries to rebuild a valid SOS marker bloc. In this case, you will be prompted to enter a value for 'padding'; this will be amount of 0x00 bytes prepended to the recovered data.