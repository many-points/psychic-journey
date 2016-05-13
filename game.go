package main

import (
  "image"
  "image/color"
  "fmt"
  "math"

  "azul3d.org/engine/gfx"
  "azul3d.org/engine/gfx/window"
  "azul3d.org/engine/gfx/gfxutil"
  "azul3d.org/engine/keyboard"

  "github.com/llgcode/draw2d"
  "github.com/llgcode/draw2d/draw2dimg"
)

type p struct {

  x   float64
  y   float64
  dx  float64
  dy  float64

  a   bool

}

func makep() *p {
  p := &p{
    x:  100,
    y:  400,
    dx: 0,
    dy: 0,
    a:  false,
  }
  return p
}

func createCard(p *p, d gfx.Device) *gfx.Object {
  cardMesh := gfx.NewMesh()
  cardMesh.Vertices = []gfx.Vec3{
    // Left triangle
    {-1, 1, 0},  // Left-Top
    {-1, -1, 0}, // Left-Bottom
    {1, -1, 0},  // Right-Bottom
    // Right triangle.
    {-1, 1, 0}, // Left-Top
    {1, -1, 0}, // Right-Bottom
    {1, 1, 0},  // Right-Top
  }
  cardMesh.TexCoords = []gfx.TexCoordSet{
    {
      Slice: []gfx.TexCoord{
        // Left triangle.
        {0, 0},
        {0, 1},
        {1, 1},
        // Right triangle.
        {0, 0},
        {1, 1},
        {1, 0},
      },
    },
  }

  shader, err := gfxutil.OpenShader("shader")
  if err != nil {
    fmt.Println(err)
  }

  tex := gfx.NewTexture()
  tex.MinFilter = gfx.Nearest
  tex.MagFilter = gfx.Nearest
  tex.KeepDataOnLoad = true
  img := draw(p, d.Bounds())
  tex.Source = img

  card := gfx.NewObject()
  card.State = gfx.NewState()
  card.Shader = shader
  card.Textures = []*gfx.Texture{tex}
  card.Meshes = []*gfx.Mesh{cardMesh}

  return card
}

func draw(p *p, bounds image.Rectangle) *image.RGBA {
  img := image.NewRGBA(bounds)

  gc := draw2dimg.NewGraphicContext(img)

  gc.SetFillColor(color.RGBA{120,120,120,255})
  gc.SetStrokeColor(color.RGBA{120,120,120,255})
  gc.SetLineWidth(0)

  xbound := float64(bounds.Max.X)
  ybound := float64(bounds.Max.Y)

  gc.MoveTo(0, ybound-50)
  gc.LineTo(xbound, ybound-50)
  gc.LineTo(xbound, ybound)
  gc.LineTo(0, ybound)
  gc.Close()
  gc.FillStroke()

  gc.SetFillColor(color.White)
  gc.SetStrokeColor(color.White)

  gc.MoveTo(p.x, p.y)
  gc.LineTo(p.x-20, p.y)
  gc.LineTo(p.x-20, p.y-20)
  gc.LineTo(p.x, p.y-20)
  gc.Close()
  gc.FillStroke()

  if p.x < 20 {
    gc.MoveTo(xbound+p.x, p.y)
    gc.LineTo(xbound+p.x-20, p.y)
    gc.LineTo(xbound+p.x-20, p.y-20)
    gc.LineTo(xbound+p.x, p.y-20)
    gc.Close()
    gc.FillStroke()
  }

  gc.SetLineWidth(1)

  startx := 125.
  starty := 100.

  dx := p.dx * 7
  dy := p.dy * 7

  cos := 0.144
  sin := 0.083

  x1 := startx + dx - (dx * cos + dy * -sin)
  y1 := starty + dy - (dx * sin + dy * cos)

  x2 := startx + dx - (dx * cos + dy * sin)
  y2 := starty + dy - (dx * -sin + dy * cos)

  gc.MoveTo(startx, starty)
  gc.LineTo(startx + dx, starty + dy)
  gc.Close()
  gc.FillStroke()

  gc.MoveTo(startx + dx, starty + dy)
  gc.LineTo(x1, y1)
  gc.LineTo(x2, y2)
  gc.LineTo(startx + dx, starty + dy)
  gc.Close()
  gc.FillStroke()

  draw2d.SetFontFolder(".")
  gc.SetFontData(draw2d.FontData{Name: "courier", Family: 0, Style: draw2d.FontStyleBold})
  gc.SetFontSize(10)
  gc.FillStringAt(fmt.Sprintf("%6.2f %5.2f %5.2f %5.2f %v", p.x, p.y, p.dx, p.dy, p.a), 10, 10)
  gc.FillStringAt(fmt.Sprintf("%6.2f", math.Sqrt(p.dx*p.dx + p.dy*p.dy)), 10, 20)

  return img
}

func gfxLoop(w window.Window, d gfx.Device) {
  p := makep()

  card := createCard(p, d)

  events := make(chan window.Event, 64)
  w.Notify(events, window.KeyboardButtonEvents)

  watcher := keyboard.NewWatcher()

  for {
    window.Poll(events, func(e window.Event) {
      switch ev := e.(type) {
      case keyboard.ButtonEvent:
        watcher.SetState(ev.Key, ev.State)
      }
    })

    //Inertia after stopping
    if watcher.Up(keyboard.ArrowLeft) && watcher.Up(keyboard.ArrowRight) {
      delta := math.Log(math.Abs(p.dx)+1)/4
      if p.dx > 0 {
        p.dx -= delta
      } else
      if p.dx < 0 {
        p.dx += delta
      }
    } else

    if watcher.Down(keyboard.ArrowLeft) && watcher.Down(keyboard.ArrowRight) {
      delta := math.Log(math.Abs(p.dx)+1)/4
      if p.dx > 0 {
        p.dx -= delta
      } else
      if p.dx < 0 {
        p.dx += delta
      }
    } else

    //Accelerate left
    if watcher.Down(keyboard.ArrowLeft) && watcher.Up(keyboard.ArrowRight) {
      p.dx -= 1
    } else

    //Accelerate right
    if watcher.Down(keyboard.ArrowRight) && watcher.Up(keyboard.ArrowLeft) {
      p.dx += 1
    }

    //Jumping impulse
    if watcher.Down(keyboard.ArrowUp) && !p.a {
      p.dy -= 7
    }

    //If velocity is low enough, it becomes zero
    if math.Abs(p.dx) < 0.1 {
      p.dx = 0
    }

    if math.Abs(p.dy) < 0.1 {
      p.dy = 0
    }

    //Maximum x velocity
    if math.Abs(p.dx) > 7 {
      if p.dx > 0 {
        p.dx = 7
      } else
      if p.dx < 0 {
        p.dx = -7
      }
    }

    //Move position
    p.x += p.dx
    p.dy += 0.3
    p.y += p.dy

    ybound := float64(d.Bounds().Max.Y)
    xbound := float64(d.Bounds().Max.X)

    //On the ground
    if p.y >= ybound - 50 {
      p.y = ybound - 50
      p.dy = 0
      p.a = false
    } else

    //in the air
    if p.y < ybound - 50 {
      p.a = true
    }

    //Screen loop
    if p.x < 0 {
      p.x = xbound
    }

    if p.x > xbound {
      p.x = 0
    }

//  d.Clear(d.Bounds(), gfx.Color{1, 1, 1, 1})
    d.ClearDepth(d.Bounds(), 1.0)

    img := draw(p, d.Bounds())
    card.Textures[0].Loaded = false
    card.Textures[0].Source = img

    d.Draw(d.Bounds(), card, nil)

    d.Render()
  }
}

func main() {
  window.Run(gfxLoop, nil)
}
