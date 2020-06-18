package messages

import (
	"errors"
	"io"

	"github.com/clearmatics/autonity/rlp"
)

type Validatable interface {
	Valid() bool
}

type NodeId interface {
}

type Proposal struct {
	Round      int64
	Height     uint64
	ValidRound int64
	NodeId     NodeId
	Value      Validatable
}

type Prevote struct {
	Round     int64
	Height    uint64
	NodeId    NodeId
	ValueHash []byte
}

type Precommit struct {
	Prevote
}

// RLP encoding doesn't support negative big.Int, so we have to pass one additionnal field to represents validRound = -1.
// Note that we could have as well indexed rounds starting by 1, but we want to stay close as possible to the spec.
func (p *Proposal) EncodeRLP(w io.Writer) error {
	isValidRoundNil := false
	var validRound uint64
	if p.ValidRound == -1 {
		validRound = 0
		isValidRoundNil = true
	} else {
		validRound = uint64(p.ValidRound)
	}

	return rlp.Encode(w, []interface{}{
		uint64(p.Round),
		p.Height,
		validRound,
		isValidRoundNil,
		p.Value,
	})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (p *Proposal) DecodeRLP(s *rlp.Stream) error {
	var proposal struct {
		Round           uint64
		Height          uint64
		ValidRound      uint64
		IsValidRoundNil bool
		Value           Validatable
	}

	if err := s.Decode(&proposal); err != nil {
		return err
	}
	var validRound int64
	if proposal.IsValidRoundNil {
		if proposal.ValidRound != 0 {
			return errors.New("bad proposal validRound with isValidround nil")
		}
		validRound = -1
	} else {
		validRound = int64(proposal.ValidRound)
	}

	if proposal.Value == nil {
		return errors.New("bad proposal with nil value")
	}

	p.Round = int64(proposal.Round)
	p.Height = proposal.Height
	p.ValidRound = validRound
	p.Value = proposal.Value

	return nil
}

// // EncodeRLP serializes b into the Ethereum RLP format.
// func (sub *Vote) EncodeRLP(w io.Writer) error {
// 	return rlp.Encode(w, []interface{}{uint64(sub.Round), sub.Height, sub.ProposedBlockHash})
// }

// // DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
// func (sub *Vote) DecodeRLP(s *rlp.Stream) error {
// 	var vote struct {
// 		Round             uint64
// 		Height            *big.Int
// 		ProposedBlockHash common.Hash
// 	}

// 	if err := s.Decode(&vote); err != nil {
// 		return err
// 	}
// 	sub.Round = int64(vote.Round)
// 	if sub.Round > MaxRound {
// 		return errInvalidMessage
// 	}
// 	sub.Height = vote.Height
// 	sub.ProposedBlockHash = vote.ProposedBlockHash
// 	return nil
// }

// func (sub *Vote) String() string {
// 	return fmt.Sprintf("{Round: %v, Height: %v ProposedBlockHash: %v}", sub.Round, sub.Height, sub.ProposedBlockHash.String())
// }
