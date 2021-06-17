package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// ContractID is the ID for a Hedera smart contract
type ContractID struct {
	Shard    uint64
	Realm    uint64
	Contract uint64
}

// ContractIDFromString constructs a ContractID from a string formatted as `Shard.Realm.Contract` (for example "0.0.3")
func ContractIDFromString(data string) (ContractID, error) {
	checksum, err := checksumParseAddress("", data)
	if err != nil {
		return ContractID{}, err
	}

	err = checksumVerify(checksum.status)
	if err != nil {
		return ContractID{}, err
	}

	return ContractID{
		Shard:    uint64(checksum.num1),
		Realm:    uint64(checksum.num2),
		Contract: uint64(checksum.num3),
	}, nil
}

// ContractIDFromSolidityAddress constructs a ContractID from a string representation of a solidity address
func ContractIDFromSolidityAddress(s string) (ContractID, error) {
	shard, realm, contract, err := idFromSolidityAddress(s)
	if err != nil {
		return ContractID{}, err
	}

	return ContractID{
		Shard:    shard,
		Realm:    realm,
		Contract: contract,
	}, nil
}

// String returns the string representation of a ContractID formatted as `Shard.Realm.Contract` (for example "0.0.3")
func (id ContractID) String() string {
	checksum, err := checksumParseAddress("", fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract))
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", id.Shard, id.Realm, id.Contract)
	}
	return fmt.Sprintf("%d.%d.%d-%s", id.Shard, id.Realm, id.Contract, checksum.correctChecksum)
}

// ToSolidityAddress returns the string representation of the ContractID as a solidity address.
func (id ContractID) ToSolidityAddress() string {
	return idToSolidityAddress(id.Shard, id.Realm, id.Contract)
}

func (id ContractID) toProtobuf() *proto.ContractID {
	return &proto.ContractID{
		ShardNum:    int64(id.Shard),
		RealmNum:    int64(id.Realm),
		ContractNum: int64(id.Contract),
	}
}

func contractIDFromProtobuf(pb *proto.ContractID) ContractID {
	if pb == nil {
		return ContractID{}
	}
	return ContractID{
		Shard:    uint64(pb.ShardNum),
		Realm:    uint64(pb.RealmNum),
		Contract: uint64(pb.ContractNum),
	}
}

func (id ContractID) toProtoKey() *proto.Key {
	return &proto.Key{Key: &proto.Key_ContractID{ContractID: id.toProtobuf()}}
}

func (id ContractID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func ContractIDFromBytes(data []byte) (ContractID, error) {
	pb := proto.ContractID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ContractID{}, err
	}

	return contractIDFromProtobuf(&pb), nil
}
