package img

import "image"

type Img struct {
	image.Image
}

func (this *Img) Save(filename string) error {
	return Save(filename, this.Image)
}

func (this *Img) Resize(maxSize uint) {
	this.Image = Resize(this.Image, maxSize)
}

func (this *Img) DrawImg(img image.Image, offset image.Point) (err error) {
	this.Image, err = DrawImg(this.Image, img, offset)
	return err
}

//func (this *Img) Write(text string) {
//	this.Image = Resize(maxSize, this.Image)
//}
