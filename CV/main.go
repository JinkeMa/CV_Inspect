package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gocv.io/x/gocv"
)

type ROI_AnL struct {
	image     gocv.Mat
	imageH    gocv.Mat
	imageS	  gocv.Mat
	imageV 	  gocv.Mat
	target	  gocv.Mat
	roiImageU gocv.Mat
	roi       image.Rectangle
	pixel 	  int

}

var root = flag.String("R", "", "Test image set path")
var sample_name = flag.String("S", "", "Sample image path for roi selection")
var BatchNum = flag.String("P", "1", "Wrapping analysis result number")

//MemCon plays the role for filtering the reflection based on the memory convolution kernel method---created by asasi:)
func (roi *ROI_AnL) MemCon(size int)gocv.Mat{
	temp:=roi.roiImageU.Region(image.Rect(0, 0, size, roi.roi.Dy()))
	kernel:=temp.Clone()
	for i := 1; i < int(roi.roi.Dx()/size)-2; i++{
		gocv.Multiply(kernel,roi.roiImageU.Region(image.Rect(size*i, 0, size*(i+1), roi.roi.Dy())), &kernel)
	}

	return kernel
}
//SavingToExcel would save the analysis result into the excel sheet
func (roi *ROI_AnL) SavingToExcel(seq int,path string, BatchNum string)error{
	var file *excelize.File

	if seq==1{
		file=excelize.NewFile()
		file.NewSheet("Wrapping analysis result")
	}else{
		file,_= excelize.OpenFile(fmt.Sprintf("Analysis_Result%s.xlsx", BatchNum))
	}
	if err:=file.SetCellValue("Wrapping analysis result", fmt.Sprintf("A%d", seq), filepath.Base(path));err!=nil{
		log.Println("Set cell value error:",err)
	}
	log.Println(fmt.Sprintf("C%d", seq),path)
	if err:=file.AddPicture("Wrapping analysis result",fmt.Sprintf("C%d", seq),path, `{"x_scale": 0.1, "y_scale": 0.1}`);err!=nil{
		log.Println("Insert image error:",err)
	}
	if err:=file.SaveAs(fmt.Sprintf("Analysis_Result%s.xlsx", BatchNum)); err!=nil{
		return err
	}
	log.Println("Saved to the Excel file")
	return nil
}


func (roi *ROI_AnL) DrawLine(index interface{}) {
	Y := roi.roi.Min.Y + index.(int)
	gocv.Line(&roi.image, image.Point{X: roi.roi.Min.X, Y: Y}, image.Point{X: roi.roi.Max.X, Y: Y}, color.RGBA{0, 255, 0, 0}, 2)
}
func (roi *ROI_AnL) SaveImg(path string) string{
	//gocv.CvtColor(roi.image, &roi.image, gocv.ColorHSVToBGR)
	gocv.IMWrite(path[:len(path)-4]+"_result.jpg", roi.image)
	return path[:len(path)-4]+"_result.jpg"
}
func (roi *ROI_AnL) ShowImage1() {
	window := gocv.NewWindow("final_result")
	//gocv.IMWrite("/Users/jindongwu/Desktop/1.png", roi.image)
	for {
		window.IMShow(roi.roiImageU)
		if window.WaitKey(1) > 0 {
			break
		}

	}
}
func (roi *ROI_AnL) ShowImage2() {
	window := gocv.NewWindow("final_result")
	//gocv.IMWrite("/Users/jindongwu/Desktop/1.png", roi.image)
	for {
		window.IMShow(roi.image)
		if window.WaitKey(1) > 0 {
			break
		}

	}
}
func ShowImg(img gocv.Mat){
	window := gocv.NewWindow("test")
	//gocv.IMWrite("/Users/jindongwu/Desktop/1.png", roi.image)
	for {
		window.IMShow(img)
		if window.WaitKey(1) > 0 {
			break
		}

	}
}

func (roi *ROI_AnL) RoiImage() {
	croppedMat := roi.target.Region(roi.roi)
	roi.roiImageU = croppedMat.Clone()
	defer croppedMat.Close()
}

func (roi *ROI_AnL) FindingEdge() (int) {
	var i int
	kernel:=gocv.GetStructuringElement(gocv.MorphCross, image.Point{5,3})
	grad:=roi.roiImageU.Clone()
	
	gocv.Erode(grad, &grad, kernel)
	gocv.Sobel(grad, &grad, gocv.MatTypeCV16S, 0,1, 5, 1, 0, gocv.BorderDefault)
	temp := gocv.NewMat()
	defer temp.Close()
	grad.ConvertTo(&temp, gocv.MatTypeCV32F)
	roi.roiImageU=grad
	ke:=roi.MemCon(10)
	//depth == CV_8U, CV_32F
	gocv.Reduce(ke, &ke, 1, gocv.ReduceAvg, gocv.MatTypeCV32F)
	dstarray, _ := ke.DataPtrFloat32()
	for i,v:=range dstarray{
		if math.Abs(float64(v))!=0{
			return i
		}
	}
return i
}


func RoiSelection(img gocv.Mat) image.Rectangle {
	window := gocv.NewWindow("x_ray2")
	roi := window.SelectROI(img)
	defer window.Close()
	return roi
}

func main() {
	const Scalar float64=0.5
	flag.Parse()
	// if *root==""||*sample_name==""{
	// 	log.Println("Please input the image path")
	// 	return

	// }
	sample := gocv.IMRead(*sample_name, gocv.IMReadColor)
	if sample.Empty(){
		log.Println("Please input valid image path!")
		return
	}


	gocv.Resize(sample, &sample, image.Point{X: 0, Y: 0}, Scalar, Scalar, gocv.InterpolationDefault)
	//ShowImg(sample)
	
	//ShowImg(sample)
	
	


	roi := RoiSelection(sample)
	//map slice for storing the image path and the opencv Mat
	images := make(map[string]gocv.Mat, 10)
	// 	//root:="/Users/jindongwu/src/github.com/northvolt/go_cv/x_ray2/x-ray2_412"
	// 	//List all of the images in the folder
	err := filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, "result.jpg")&&strings.HasSuffix(path, ".jpg"){
			images[path] = gocv.IMRead(path, gocv.IMReadColor)
			return nil
		}
		log.Println("It is not the image for analysising")
		return nil
	})
	if err != nil {
		panic(err)
	}
		//Get the ROI postion based on the sample image
	//ROI selection should be out of the the main loop

	index:=0
	for v, img := range images {
		
	 	if !img.Empty() {
			index++
			target:=gocv.NewMat()
			
			gocv.Resize(img, &img, image.Point{X: 0, Y: 0}, Scalar, Scalar, gocv.InterpolationDefault)
			temp:=img.Clone()
			gocv.CvtColor(img,&img, gocv.ColorBGRToHSV)
			HsvArr:=gocv.Split(img)
			gocv.MultiplyWithParams(HsvArr[0], HsvArr[1], &target, 5, gocv.MatTypeCV16S)
			roi_an := ROI_AnL{image: temp, imageH: HsvArr[0], imageS:HsvArr[1], imageV:HsvArr[2],target: target}
			roi_an.roi = roi
			roi_an.RoiImage()
			roi_an.DrawLine(roi_an.FindingEdge())
			roi_an.ShowImage2()
			roi_an.ShowImage1()
			roi_an.SavingToExcel(index,roi_an.SaveImg(v),*BatchNum )
			defer temp.Close()
		}
	}
}
