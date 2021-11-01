package models

import (
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sergeyglazyrindev/uadmin/core"
	"gorm.io/gorm"
	"time"
)

type ScienceCategory struct {
	core.Model
	Name        string `gorm:"not null" uadmin:"list,search"`
	Icon        string `uadminform:"ScienceCategoryPhotoFormOptions" gorm:"default:''"`
}

type ScienceCategoryLocalized struct {
	core.Model
	ScienceCategory ScienceCategory
	ScienceCategoryID uint
	LanguageCode string `gorm:"size=4" uadmin:"inline" uadminform:"RequiredSelectFieldOptions"`
	Name        string `gorm:"not null" uadmin:"inline"`
	Description string `uadminform:"TextareaFieldOptions" uadmin:"inline"`
}

func (sc *ScienceCategory) String() string {
	return sc.Name
}

type ScienceTermType uint

func (ScienceTermType) Axiom() ScienceTermType {
	return 1
}

func (ScienceTermType) Law() ScienceTermType {
	return 2
}

func (ScienceTermType) Deduction() ScienceTermType {
	return 3
}

func (ScienceTermType) Induction() ScienceTermType {
	return 4
}

func HumanizeScienceTermType(termType ScienceTermType) string {
	switch termType {
	case 4:
		return "induction"
	case 2:
		return "law"
	case 3:
		return "deduction"
	default:
		return "axiom"
	}
}

type ReasonType uint

func (ReasonType) Contradicts() ReasonType {
	return 1
}

func (ReasonType) Conforms() ReasonType {
	return 2
}

func HumanizeReasonType(reason ReasonType) string {
	switch reason {
	case 2:
		return "conforms"
	default:
		return "contradicts"
	}
}

type ScienceTerm struct {
	core.Model
	Alias             string          `uadmin:"list,search" gorm:"uniqueIndex;not null"`
	Type              ScienceTermType `uadmin:"list,search" uadminform:"SelectFieldOptions"`
	Discussion        *Discussion     `uadmin:"list" uadminform:"ForeignKeyWithAutocompleteFieldOptions"`
	DiscussionID      uint            `gorm:"default:null;uniqueIndex;" uadmin:"search"`
	ScienceCategory   ScienceCategory `uadmin:"list,search" uadminform:"ForeignKeyWithAutocompleteFieldOptions"`
	ScienceCategoryID uint
}

type ScienceTermLocalized struct {
	core.Model
	ScienceTerm ScienceTerm
	ScienceTermID uint
	LanguageCode string `gorm:"size=4" uadmin:"inline" uadminform:"RequiredSelectFieldOptions"`
	ShortDescription  string          `uadmin:"inline" uadminform:"TextareaFieldOptions"`
	Description       string          `uadmin:"inline" uadminform:"TextareaFieldOptions"`
}

func (st *ScienceTerm) String() string {
	return st.Alias
}

func (st *ScienceTerm) BeforeCreate(tx *gorm.DB) error {
	if st.ID == 0 {
		//uname, err := username.NewUCDUsername()
		//if err != nil {
		//	// log.Fatal(err)
		//}
		//uname.Debug = false
		//uname.AllowSpaces = false
		//uname.AllowPunctuation = false
		//safe, err := uname.Translate(st.ShortDescription)
		st.Alias = core.ASCIIRegex.ReplaceAllLiteralString(st.Alias, "")
		st.Alias = core.WhitespaceRegex.ReplaceAllLiteralString(st.Alias, "")
	}
	return nil
}

type Expert struct {
	core.Model
	User   core.User `uadmin:"list,search" uadminform:"FkReadonlyFieldOptions"`
	UserID uint      `gorm:"uniqueIndex"`
	Languages []core.Language `gorm:"foreignKey:ID;many2many:expert_languages;"`
}

func (e *Expert) String() string {
	return fmt.Sprintf("Expert - %s", e.User.String())
}

type ExpertScienceCategory struct {
	core.Model
	Expert            Expert          `uadmin:"inline" uadminform:"ForeignKeyFieldOptions"`
	ExpertID          uint            `gorm:"index:expert_science_category,unique"`
	ScienceCategory   ScienceCategory `uadmin:"inline" uadminform:"ForeignKeyFieldOptions"`
	ScienceCategoryID uint            `gorm:"index:expert_science_category,unique"`
	Approved          bool            `uadmin:"inline"`
}

func (esc *ExpertScienceCategory) String() string {
	return fmt.Sprintf("%s for category %s", esc.Expert.String(), esc.ScienceCategory.String())
}

type Discussion struct {
	core.Model
	Author   Expert `uadmin:"list,search" uadminform:"ForeignKeyReadonlyFieldOptions"`
	AuthorID uint
}

type DiscussionLocalized struct {
	core.Model
	Discussion Discussion
	DiscussionID uint
	LanguageCode string `gorm:"size=4" uadmin:"inline" uadminform:"RequiredSelectFieldOptions"`
	Subject  string `uadmin:"inline"`
}

func (d *Discussion) String() string {
	// d.Subject
	return fmt.Sprintf("Discussion started by author %s", d.Author.String())
}

type DiscussionComment struct {
	ID            uint       `gorm:"primarykey" nestedset:"id"`
	Discussion    Discussion `uadmin:"list,search" uadminform:"ForeignKeyReadonlyFieldOptions"`
	DiscussionID  uint       `nestedset:"scope" uadmin:"search"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt     `gorm:"index"`
	Parent        *DiscussionComment `uadminform:"ForeignKeyReadonlyFieldOptions"`
	ParentID      sql.NullInt64      `nestedset:"parent_id"`
	Rgt           int                `nestedset:"rgt"`
	Lft           int                `nestedset:"lft"`
	Depth         int                `nestedset:"depth"`
	ChildrenCount int                `nestedset:"children_count"`
	Author        Expert             `uadmin:"list,search" uadminform:"ForeignKeyReadonlyFieldOptions"`
	AuthorID      uint
}

func (d *DiscussionComment) String() string {
	// d.Subject
	return fmt.Sprintf("Discussion comment posted by author %s", d.Author.String())
}

type DiscussionCommentLocalized struct {
	core.Model
	DiscussionComment DiscussionComment
	DiscussionCommentID uint
	LanguageCode string `gorm:"size=4" uadmin:"inline" uadminform:"RequiredSelectFieldOptions"`
	Subject       string             `uadmin:"inline"`
	Body          string             `uadmin:"inline" uadminform:"TextareaFieldOptions"`
	Approved bool
}

type ExpertReview struct {
	Veracity      uint        `uadmin:"inline"`
	Reason        ReasonType  `uadmin:"inline" uadminform:"SelectFieldOptions"`
	Explanation   string      `uadmin:"inline" uadminform:"TextareaFieldOptions"`
	ScienceTerm   ScienceTerm `uadmin:"inline" uadminform:"ForeignKeyWithAutocompleteFieldOptions"`
	ScienceTermID uint
}

type DiscussionReview struct {
	core.Model
	ExpertReview
	Author       Expert `uadmin:"inline" uadminform:"ForeignKeyWithAutocompleteFieldOptions"`
	AuthorID     uint   `gorm:"index:discussion_review_author"`
	Discussion   Discussion
	DiscussionID uint `gorm:"index:discussion_review_author"`
}

type DiscussionCommentReview struct {
	core.Model
	ExpertReview
	Author              Expert `uadmin:"inline" uadminform:"ForeignKeyWithAutocompleteFieldOptions"`
	AuthorID            uint   `gorm:"index:discussion_comment_review_author"`
	DiscussionComment   DiscussionComment
	DiscussionCommentID uint `gorm:"index:discussion_comment_review_author"`
}

type ProofItApp struct {
	ID        string `sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
	Name      string         `uadmin:"list"`
	PublicKey string         `uadminform:"ReadonlyField" gorm:"uniqueIndex;not null"`
	SecretKey string         `uadminform:"ReadonlyField" gorm:"uniqueIndex;not null"`
	Languages []core.Language `gorm:"foreignKey:ID;many2many:proofitapp_languages;"`
}

func (pia *ProofItApp) String() string {
	return pia.Name
}

func (pia *ProofItApp) BeforeCreate(tx *gorm.DB) error {
	if pia.ID == "" {
		id := uuid.NewV4()
		pia.ID = id.String()
		publicKey := uuid.NewV4()
		pia.PublicKey = publicKey.String()
		secretKey := uuid.NewV4()
		pia.SecretKey = secretKey.String()
	}
	return nil
}
