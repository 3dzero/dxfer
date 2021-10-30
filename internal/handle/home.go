package handle

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yofu/dxf"
	"io"
	"net/http"
	"strconv"
)

type ITemplate interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}

func Home(T ITemplate) func(w http.ResponseWriter, _ *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var Error string
		q:=r.URL.Query()
		//filename
		filename := q.Get("filename")
		x:=q.Get("X")
		y:=q.Get("Y")
		logrus.
			WithField("x",x).
			WithField("y",y).
			WithField("filename",filename).
			Info("запрос")

		p, err:= NewPoligon(x,y)
		if err != nil{
			Error = err.Error()
		}else{
			w.Header().Set("Content-Disposition", "attachment; filename="+filename+x+"*"+y+".dxf")
			NewDXF(filename, p, w)
			return
		}

		err = T.ExecuteTemplate(w, "home", map[string]string{
			"Error": Error,
			"Filename": filename,
			"X": x,
			"Y": y,
		})
		if err != nil {
			logrus.Error(err)
		}
	}
}

type Poligon struct{
	H float64
	W float64
}

func NewDXF(filename string, poligon *Poligon, w io.Writer)error{

	d := dxf.NewDrawing()
	d.Header().LtScale = 100.0
	d.AddLayer("sq", dxf.DefaultColor, dxf.DefaultLineType, true)

	d.Line(0, 0, 0,  poligon.W, 0,0)
	d.Line(poligon.W, 0, 0,  poligon.W, poligon.H,0)

	d.Line(0, 0, 0, 0, poligon.H, 0)
	d.Line(0, poligon.H, 0, poligon.W, poligon.H, 0)
	defer d.Close()

	_,err := d.WriteTo(w)
	return err
}


func NewPoligon(x string, y string)(p *Poligon, err error){
	W, err := strconv.ParseFloat(x, 64)
	if err !=nil{
		logrus.WithField("x", x).Error(err)
		err= fmt.Errorf("ошибка размера по X")
		return
	}

	H, err := strconv.ParseFloat(y, 64)
	if err !=nil{
		logrus.WithField("y", y).Error(err)
		err= fmt.Errorf("ошибка размера по Y ")
		return
	}
	return &Poligon{H,W},nil
}
