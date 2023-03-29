package math1_v1

import (
	context "context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var e error

type Addition struct {
	gorm.Model
	A     int64
	B     int64
	Value int64
}

type Subtraction struct {
	gorm.Model
	A     int64
	B     int64
	Value int64
}

type Multiplication struct {
	gorm.Model
	A     int64
	B     int64
	Value int64
}

type Division struct {
	gorm.Model
	A     int64
	B     int64
	Value int64
}

type SvcConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"dbName"`
	DbUser   string `yaml:"dbUser"`
	Password string `yaml:"password"`
}

func getConfig() *SvcConfig {
	var cfg SvcConfig
	err := ReadFile("CONFIG_PATH", &cfg)
	if err != nil {
		log.Fatalf("error reading config")
	}
	return &cfg
}

func ReadFile(filePath string, cfg interface{}) error {
	path, found := os.LookupEnv(filePath)
	if !found {
		return errors.New("config file not found")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	ReadEnv(cfg)
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func ReadEnv(cfg interface{}) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	DatabaseConnection()
}

func DatabaseConnection() {

	cfg := getConfig()
	fmt.Println("error")
	fmt.Println(cfg.DbUser)
	dsn := cfg.DbUser + ":" + cfg.Password + "@tcp" + "(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.DbName + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn : ", dsn)

	DB, e = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(Addition{})
	DB.AutoMigrate(Subtraction{})
	DB.AutoMigrate(Multiplication{})
	DB.AutoMigrate(Division{})

	if e != nil {
		log.Fatal("Error connecting to the database...", e)
	}
	fmt.Println("Database connection successful...", e)
}

type Server struct {
	UnimplementedMathServiceServer
}

func (s *Server) Add(ctx context.Context, req *Request) (*Response, error) {
	a, b := req.GetNum1(), req.GetNum2()
	value := int64(a + b)
	data := Addition{
		A:     int64(a),
		B:     int64(b),
		Value: value,
	}
	fmt.Println(DB)
	rows := DB.Create(&data)
	if rows == nil {
		fmt.Println("Error in row creation", e)
		return nil, e
	}
	return &Response{Result: value}, nil
}

func (s *Server) Subtract(ctx context.Context, req *Request) (*Response, error) {
	a, b := req.GetNum1(), req.GetNum2()
	value := int64(a - b)
	data := Subtraction{
		A:     int64(a),
		B:     int64(b),
		Value: value,
	}
	fmt.Println(DB)
	rows := DB.Create(&data)
	if rows == nil {
		fmt.Println("Error in row creation", e)
		return nil, e
	}
	return &Response{Result: value}, nil
}

func (s *Server) Multiply(ctx context.Context, req *Request) (*Response, error) {
	a, b := req.GetNum1(), req.GetNum2()
	value := int64(a * b)
	data := Multiplication{
		A:     int64(a),
		B:     int64(b),
		Value: value,
	}
	fmt.Println(DB)
	rows := DB.Create(&data)
	if rows == nil {
		fmt.Println("Error in row creation", e)
		return nil, e
	}
	return &Response{Result: value}, nil
}

func (s *Server) Divide(ctx context.Context, req *Request) (*Response, error) {
	a, b := req.GetNum1(), req.GetNum2()
	// fmt.Println(b)
	if b == 0 {
		fmt.Println("Denominator:", b)
		fmt.Println("Division not possible")
		return &Response{}, nil
	} else {
		value := int64(a / b)
		fmt.Println(value)
		data := Division{
			A:     int64(a),
			B:     int64(b),
			Value: value,
		}
		rows := DB.Create(&data)
		if rows == nil {
			fmt.Println("Error in row creation", e)
			return nil, e
		}
		return &Response{Result: value}, nil
	}
}
