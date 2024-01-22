package storage

import (
	"bufio"
	"encoding/json"
	"os"
)

type fileStorage struct {
	consumer     *Consumer
	producer     *Producer
	inMemoryData map[string]*ShortURL
}

func (f *fileStorage) Find(uuid string) (*ShortURL, bool) {
	for {
		event, err := f.consumer.ReadItem()

		if err != nil {
			break
		}

		f.inMemoryData[event.UUID] = event
	}

	item, ok := f.inMemoryData[uuid]

	if ok {
		return item, true
	}

	return nil, false
}

func (f *fileStorage) Save(url *ShortURL) error {
	return f.producer.WriteEvent(url)
}

func (f *fileStorage) Size() int {
	return len(f.inMemoryData)
}

func (f *fileStorage) Ping() bool {
	return true
}

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

func (p *Producer) WriteEvent(event *ShortURL) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{file: file, writer: bufio.NewWriter(file)}, nil
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{file: file, decoder: json.NewDecoder(file)}, nil
}

func (c *Consumer) ReadItem() (*ShortURL, error) {
	event := &ShortURL{}
	if err := c.decoder.Decode(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Consumer) Done() {
	c.decoder.InputOffset()
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func NewFileStorage(filePath string) (Storage, error) {
	Producer, err := NewProducer(filePath)
	if err != nil {
		return nil, err
	}

	Consumer, err := NewConsumer(filePath)
	if err != nil {
		return nil, err
	}

	return &fileStorage{Consumer, Producer, make(map[string]*ShortURL)}, nil
}
