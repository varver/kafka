package kafka

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
)

type Request interface {
	Bytes() ([]byte, error)
	WriteTo(io.Writer) (int64, error)
}

var _ Request = &MetadataReq{}
var _ Request = &ProduceReq{}
var _ Request = &FetchReq{}
var _ Request = &ConsumerMetadataReq{}
var _ Request = &OffsetReq{}
var _ Request = &OffsetCommitReq{}
var _ Request = &OffsetFetchReq{}

func testRequestSerialization(t *testing.T, r Request) {
	var buf bytes.Buffer
	if n, err := r.WriteTo(&buf); err != nil {
		t.Fatalf("could not write request to buffer: %s", err)
	} else if n != int64(buf.Len()) {
		t.Fatalf("writer returned invalid number of bytes written %d != %d", n, buf.Len())
	}
	b, err := r.Bytes()
	if err != nil {
		t.Fatalf("could not convert request to bytes: %s", err)
	}
	if !bytes.Equal(b, buf.Bytes()) {
		t.Fatal("Bytes() and WriteTo() serialized request is of different form")
	}
}

func TestMetadataRequest(t *testing.T) {
	req1 := &MetadataReq{
		CorrelationID: 123,
		ClientID:      "testcli",
		Topics:        nil,
	}
	testRequestSerialization(t, req1)
	b, _ := req1.Bytes()
	expected := []byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x7, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6c, 0x69, 0x0, 0x0, 0x0, 0x0}

	if !bytes.Equal(b, expected) {
		t.Fatalf("expected different bytes representation: %v", b)
	}

	req2 := &MetadataReq{
		CorrelationID: 123,
		ClientID:      "testcli",
		Topics:        []string{"foo", "bar"},
	}
	testRequestSerialization(t, req2)
	b, _ = req2.Bytes()
	expected = []byte{0x0, 0x0, 0x0, 0x1f, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x7, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6c, 0x69, 0x0, 0x0, 0x0, 0x2, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x3, 0x62, 0x61, 0x72}

	if !bytes.Equal(b, expected) {
		t.Fatalf("expected different bytes representation: %v", b)
	}
}

func TestMetadataResponse(t *testing.T) {
	msgb := []byte{0x0, 0x0, 0x1, 0xc7, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0xc0, 0x10, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x12, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x11, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12}
	resp, err := ReadMetadataResp(bytes.NewBuffer(msgb))
	if err != nil {
		t.Fatalf("could not read metadata response: %s", err)
	}
	expected := &MetadataResp{
		CorrelationID: 123,
		Brokers: []BrokerMetadata{
			BrokerMetadata{NodeID: 49168, Host: "172.17.42.1", Port: 49168},
			BrokerMetadata{NodeID: 49170, Host: "172.17.42.1", Port: 49170},
			BrokerMetadata{NodeID: 49169, Host: "172.17.42.1", Port: 49169},
			BrokerMetadata{NodeID: 49171, Host: "172.17.42.1", Port: 49171},
		},
		Topics: []TopicMetadata{
			TopicMetadata{
				Name: "foo",
				Err:  error(nil),
				Partitions: []PartitionMetadata{
					PartitionMetadata{Err: error(nil), ID: 2, Leader: 49171, Replicas: []int32{49171, 49168, 49169}, Isrs: []int32{49171, 49168, 49169}},
					PartitionMetadata{Err: error(nil), ID: 5, Leader: 49170, Replicas: []int32{49170, 49168, 49169}, Isrs: []int32{49170, 49168, 49169}},
					PartitionMetadata{Err: error(nil), ID: 4, Leader: 49169, Replicas: []int32{49169, 49171, 49168}, Isrs: []int32{49169, 49171, 49168}},
					PartitionMetadata{Err: error(nil), ID: 1, Leader: 49170, Replicas: []int32{49170, 49171, 49168}, Isrs: []int32{49170, 49171, 49168}},
					PartitionMetadata{Err: error(nil), ID: 3, Leader: 49168, Replicas: []int32{49168, 49169, 49170}, Isrs: []int32{49168, 49169, 49170}},
					PartitionMetadata{Err: error(nil), ID: 0, Leader: 49169, Replicas: []int32{49169, 49170, 49171}, Isrs: []int32{49169, 49170, 49171}},
				},
			},
			TopicMetadata{
				Name: "test",
				Err:  error(nil),
				Partitions: []PartitionMetadata{
					PartitionMetadata{Err: error(nil), ID: 1, Leader: 49169, Replicas: []int32{49169, 49170, 49171}, Isrs: []int32{49169, 49170, 49171}},
					PartitionMetadata{Err: error(nil), ID: 0, Leader: 49168, Replicas: []int32{49168, 49169, 49170}, Isrs: []int32{49168, 49169, 49170}},
				},
			},
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("expected different message: %#v", resp)
	}
}

func TestProduceRequest(t *testing.T) {
	req := &ProduceReq{
		CorrelationID: 241,
		ClientID:      "test",
		RequiredAcks:  RequiredAcksAll,
		Timeout:       time.Second,
		Topics: []ProduceReqTopic{
			ProduceReqTopic{
				Name: "foo",
				Partitions: []ProduceReqPartition{
					ProduceReqPartition{
						Partition: 0,
						Messages: []*Message{
							&Message{
								Offset: 53,
								Crc:    92,
								Key:    []byte("foo"),
								Value:  []byte("bar"),
							},
						},
					},
				},
			},
		},
	}
	testRequestSerialization(t, req)
	b, _ := req.Bytes()
	expected := []byte{0x0, 0x0, 0x0, 0x49, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xff, 0xff, 0x0, 0x0, 0x3, 0xe8, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72}

	if !bytes.Equal(b, expected) {
		t.Fatalf("expected different bytes representation: %#v", b)
	}
}

func TestProduceResponse(t *testing.T) {
	msgb1 := []byte{0x0, 0x0, 0x0, 0x22, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x6, 0x66, 0x72, 0x75, 0x69, 0x74, 0x73, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x5d, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	resp1, err := ReadProduceResp(bytes.NewBuffer(msgb1))
	if err != nil {
		t.Fatalf("could not read metadata response: %s", err)
	}
	expected1 := &ProduceResp{
		CorrelationID: 241,
		Topics: []ProduceRespTopic{
			ProduceRespTopic{
				Name: "fruits",
				Partitions: []ProduceRespPartition{
					ProduceRespPartition{
						Partition: 93,
						Err:       ErrUnknownTopicOrPartition,
						Offset:    -1,
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp1, expected1) {
		t.Fatalf("expected different message: %#v", resp1)
	}

	msgb2 := []byte{0x0, 0x0, 0x0, 0x1f, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}
	resp2, err := ReadProduceResp(bytes.NewBuffer(msgb2))
	if err != nil {
		t.Fatalf("could not read metadata response: %s", err)
	}
	expected2 := &ProduceResp{
		CorrelationID: 241,
		Topics: []ProduceRespTopic{
			ProduceRespTopic{
				Name: "foo",
				Partitions: []ProduceRespPartition{
					ProduceRespPartition{
						Partition: 0,
						Err:       error(nil),
						Offset:    1,
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp2, expected2) {
		t.Fatalf("expected different message: %#v", resp2)
	}
}

func TestFetchRequest(t *testing.T) {
	req := &FetchReq{
		CorrelationID: 241,
		ClientID:      "test",
		MaxWaitTime:   time.Second * 2,
		MinBytes:      12454,
		Sources: []FetchReqTopic{
			FetchReqTopic{
				Topic: "foo",
				Partitions: []FetchReqPartition{
					FetchReqPartition{Partition: 421, FetchOffset: 529, MaxBytes: 4921},
					FetchReqPartition{Partition: 0, FetchOffset: 11, MaxBytes: 92},
				},
			},
		},
	}
	testRequestSerialization(t, req)
	b, _ := req.Bytes()
	expected := []byte{0x0, 0x0, 0x0, 0x47, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x7, 0xd0, 0x0, 0x0, 0x30, 0xa6, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x1, 0xa5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x11, 0x0, 0x0, 0x13, 0x39, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xb, 0x0, 0x0, 0x0, 0x5c}

	if !bytes.Equal(b, expected) {
		t.Fatalf("expected different bytes representation: %#v", b)
	}
}

func TestFetchResponse(t *testing.T) {
	msgb := []byte{0x0, 0x0, 0x0, 0x75, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0}
	resp, err := ReadFetchResp(bytes.NewBuffer(msgb))
	if err != nil {
		t.Fatalf("could not read metadata response: %s", err)
	}
	expected := &FetchResp{
		CorrelationID: 241,
		Sources: []FetchRespTopic{
			FetchRespTopic{
				Topic: "foo",
				Partitions: []FetchRespPartition{
					FetchRespPartition{
						Partition: 0,
						Err:       error(nil),
						TipOffset: 4,
						Messages: []*Message{
							&Message{Offset: 2, Crc: 0xb8ba5f57, Key: []uint8{0x66, 0x6f, 0x6f}, Value: []uint8{0x62, 0x61, 0x72}},
							&Message{Offset: 3, Crc: 0xb8ba5f57, Key: []uint8{0x66, 0x6f, 0x6f}, Value: []uint8{0x62, 0x61, 0x72}},
						},
					},
					FetchRespPartition{
						Partition: 1,
						Err:       ErrUnknownTopicOrPartition,
						TipOffset: -1,
						Messages:  []*Message{},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("expected different message: %#v", resp)
	}
}
