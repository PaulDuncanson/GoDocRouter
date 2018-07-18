# GoDocRouter
GoDocRouter is an application framework written in Golang that routes documents and images through a configurable workflow process.<p>
This workflow defines the order of events by utilizing the Open Source designer tool, <a href=https://en.wikipedia.org/wiki/LibreOffice>LibreOffice</a>, with corresponding connected shapes.<p>
<p align="center">
  <img src="./GoDocRouter.png" width="500"/>
</p>
One transformation defined in the workflow as illustrated above, doOCRToText, will perform Optical Character Recognition on supplied document images utilizing the <b><a href=https://en.wikipedia.org/wiki/Tesseract_(software)>Tesseract-OCR</a></b> library.<p><p>
<b>To install and run:</b><p><p>
<b>Step 1:</b>  Download and unzip the attached GoDocRouter Application by clicking the corresponding 'Clone or Download...' green button and then select 'Download Zip'.<p>
<b>Step 2:</b>  Install the necessary Golang compiler that you can find <b><a href=https://golang.org/dl/>here.</a></b><p>
<b>Step 3:</b>  Open a terminal (Linux/Mac) or Command Line (Windows) and change your current directory where the files have been unzipped.<p>
<b>Step 4:</b>  Type the following to compile and run: >go run CADEngine.go<p>
